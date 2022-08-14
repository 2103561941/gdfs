package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/cyb0225/gdfs/proto/namenode"
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
		log.Fatalf("failed to open file: %s: %s", localFilePath, err.Error())
	}
	fileinfo, err := fd.Stat()
	if err != nil {
		log.Fatalf("failed to get file: %s stat: %s", localFilePath, err.Error())
	}
	// get bytes, should transform to KB
	filesize := (float64(fileinfo.Size()) / 1024) 
	fmt.Println("filesize: ", filesize)

	res, err := put(remoteFilePath, filesize)
	if err != nil {
		log.Fatalf("get datanode information from namenode failed: %s\n", err.Error())
	}

	fmt.Printf("client get server: %+v", res)
}

func put(filepath string, filesize float64) (*pb.PutResponse, error){
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to server failed: %w", err)
	}
	defer conn.Close()

	c := pb.NewNameNodeClient(conn)
	req := &pb.PutRequest{
		RemoteFilePath: filepath,
		Filesize: filesize,
	}
	res, err := c.Put(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("get from %s server failed: %w", addr, err)
	}

	return res, nil
}
