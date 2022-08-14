package cmd

import (
	"context"
	"fmt"
	"log"

	pb "github.com/cyb0225/gdfs/proto/namenode"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = "127.0.0.1:50051"
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
	// localFilePath := args[0] // save remoteFile to this file
	remoteFilePath := args[1]

	res, err := get(remoteFilePath)
	if err != nil {
		log.Fatalf("get datanode information from namenode failed: %s\n", err.Error())
	}

	fmt.Printf("client get server: %+v", res)
}

// get filepath's datanodes information from namenode
func get(filepath string) (*pb.GetResponse, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to namenode[%s] failed: %w", addr, err)
	}

	defer conn.Close()

	c := pb.NewNameNodeClient(conn)
	res, err := c.Get(context.Background(), &pb.GetRequest{RemoteFilePath: filepath})
	if err != nil {
		return nil, fmt.Errorf("get datanodes' information failed: %w", err)
	}

	return res, nil
}
