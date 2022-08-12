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
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect to server failed: %s", err.Error())
	}
	defer conn.Close()

	c := pb.NewNameNodeClient(conn)
	res, err := c.Put(context.Background(), &pb.PutRequest{RemoteFilePath: "/tmp/cyb"})
	if err != nil {
		log.Fatalf("get from %s server failed: %s", addr, err.Error())
	}

	fmt.Printf("client get server: %+v", res)
}
