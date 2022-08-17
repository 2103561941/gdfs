package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cyb0225/gdfs/internal/client/config"
	"github.com/cyb0225/gdfs/pkg/log"
	pb2 "github.com/cyb0225/gdfs/proto/datanode"
	pb1 "github.com/cyb0225/gdfs/proto/namenode"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get file from remote file system and save in local file system",

	// args check
	Args: func(cmd *cobra.Command, args []string) error {
		return equalNumArgs(2, args)
	},
	Run: Get,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func Get(cmd *cobra.Command, args []string) {
	localFilePath := args[0] // save remoteFile to this file
	remoteFilePath := args[1]

	res, err := get(remoteFilePath)
	if err != nil {
		log.Fatal("connect to namenode failed", log.Err(err))
	}

	fd, err := os.OpenFile(localFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("create file failed", log.String("file", localFilePath), log.Err(err))
	}
	defer fd.Close()

	// read the chunked file and flatten the complete content
	w := bufio.NewWriter(fd)
	chunks := res.Chunks
	for i := 0; i < len(chunks); i++ {
		filekey := chunks[i].FileKey
		backups := chunks[i].Backups
		isError := true
		for j := 0; j < len(backups); j++ {
			if err := getdata(filekey, backups[j], w); err != nil {
				log.Error("client get file failed", log.String("datanod", backups[j]), log.Err(err))
				continue
			}
			isError = false
			break
		}
		// it means that client cannot get data from any datanode
		if isError {
			log.Fatalf("get file failed from any datanode", log.String("filekey", filekey))
		}
	}

	// Notice, if don't use this funciton, file will not have the data.
	w.Flush()
	log.Info("get file success!")
}

// get filepath's datanodes information from namenode
func get(filepath string) (*pb1.GetResponse, error) {
	conn, err := grpc.Dial(config.Cfg.NamenodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	c := pb1.NewNameNodeClient(conn)
	res, err := c.Get(context.Background(), &pb1.GetRequest{RemoteFilePath: filepath})
	if err != nil {
		return nil, fmt.Errorf("get from namenode failed: %w", err)
	}
	return res, nil
}

// read file data from datanode.
func getdata(filekey string, addr string, w io.Writer) error {

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("connect to namenode failed: %w", err)
	}

	defer conn.Close()

	c := pb2.NewDataNodeClient(conn)
	req := &pb2.GetRequset{
		Filekey: filekey,
	}
	stream, err := c.Get(context.Background(), req)
	if err != nil {
		return fmt.Errorf("get stream from %s failed: %w", addr, err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("get file data from %s failed: %w", addr, err)
		}

		// log.Println(string(res.Databytes))

		if _, err := w.Write(res.Databytes); err != nil {
			return fmt.Errorf("write to local file failed: %w", err)
		}
	}

	if err := stream.CloseSend(); err != nil {
		return fmt.Errorf("close send stream failed: %w", err)
	}

	return nil
}
