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

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "put local file to remote file system",
	Args: func(cmd *cobra.Command, args []string) error {
		return equalNumArgs(2, args)
	},
	Run: Put,
}

func init() {
	rootCmd.AddCommand(putCmd)
}

func Put(cmd *cobra.Command, args []string) {
	localFilePath := args[0]
	remoteFilePath := args[1] // choose a place to store this file

	fd, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal("failed to open file", log.String("file", localFilePath), log.Err(err))
	}
	defer fd.Close()

	fileinfo, err := fd.Stat()
	if err != nil {
		log.Fatal("get file stat failed", log.String("file", localFilePath), log.Err(err))
	}
	// get bytes, should transform to KB
	filesize := (fileinfo.Size())

	res, err := put(remoteFilePath, filesize)
	if err != nil {
		log.Fatal("connect to namenode failed", log.String("namenode", config.Cfg.NamenodeAddr) ,log.Err(err))
	}

	// put file data to datanodes
	r := bufio.NewReader(fd)
	for i := 0; i < len(res.Chunks); i++ {
		filekey := res.Chunks[i].FileKey
		backups := res.Chunks[i].Backups
		isError := true
		for j := 0; j < len(backups); j++ {
			if err := putdata(filekey, r, backups[j:]); err != nil {
				log.Error("put file to datanode failed", log.String("datanode", backups[j]), log.String("filekey", filekey), log.Err(err))
				continue
			}
			// put to datanode success. then put the next chunk
			isError = false
			break
		}
		if isError {
			log.Fatalf("put file to any datanode failed", log.String("filekey", filekey))
		}
	}
	log.Info("put file success!")
}

// get datanode information from namenode
func put(filepath string, filesize int64) (*pb1.PutResponse, error) {
	conn, err := grpc.Dial(config.Cfg.NamenodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := pb1.NewNameNodeClient(conn)
	req := &pb1.PutRequest{
		RemoteFilePath: filepath,
		Filesize:       filesize,
	}
	res, err := c.Put(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("get from [%s] server failed: %w", config.Cfg.NamenodeAddr, err)
	}

	return res, nil
}

// put data to datanode
// add[0] stored the address which will be visited, and adds[1:] stored the other backups' address.
func putdata(filekey string, r io.Reader, adds []string) error {

	// if put one datanode failed, then try to put to next backups.
	// at the same time
	address := adds[0]

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("connect to server failed: %w", err)
	}
	defer conn.Close()

	c := pb2.NewDataNodeClient(conn)
	stream, err := c.Put(context.Background())
	if err != nil {
		return fmt.Errorf("get stream failed: %w", err)
	}

	// send basic information to datanode.
	if err := stream.Send(&pb2.PutRequest{Filekey: filekey, Adds: adds[1:]}); err != nil {
		return fmt.Errorf("send filestat to datanode failed: %w", err)
	}

	buf := make([]byte, 1024) //chunk size can divide it.   chunksize mod bufsize = 0
	sum := 0                  // stored the bytes that read.
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			// size of buf is 0
			break
		}
		if err != nil {
			return fmt.Errorf("read file %s failed: %w", filekey, err)
		}

		sum += n
		if sum >= int(config.Cfg.ChunkSize) { // every chunk's size
			break
		}

		if err := stream.Send(&pb2.PutRequest{Databytes: buf[:n]}); err != nil {
			return fmt.Errorf("send filestat to datanode failed: %w", err)
		}
	}

	if _, err = stream.CloseAndRecv(); err != nil {
		return fmt.Errorf("close client stream failed: %w", err)
	}

	return nil
}
