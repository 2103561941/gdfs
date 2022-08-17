package cmd

import (
	"context"
	"fmt"

	"github.com/cyb0225/gdfs/internal/client/config"
	"github.com/cyb0225/gdfs/pkg/log"
	pb1 "github.com/cyb0225/gdfs/proto/namenode"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete file in remote file system",
	Args: func(cmd *cobra.Command, args []string) error {
		return equalNumArgs(1, args) 
	},
	Run:Delete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func Delete(cmd *cobra.Command, args []string) {
	_, err := delete(args[0])
	if err != nil {
		log.Fatal("delete file failed", log.Err(err))
	}

	log.Info("delete file success!")
}

func delete(filepath string) (*pb1.DeleteResponse, error) {
	conn, err := grpc.Dial(config.Cfg.NamenodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := pb1.NewNameNodeClient(conn)
	req := &pb1.DeleteRequest{RemoteFilePath: filepath}
	res, err := c.Delete(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("get from namenode failed: %w", err)
	}

	return res, nil
}