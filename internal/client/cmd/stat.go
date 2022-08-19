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

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "get metadata of remote file system",
	Args: func(cmd *cobra.Command, args []string) error {
		return equalNumArgs(1, args)
	},
	Run: Stat,
}

func init() {
	rootCmd.AddCommand(statCmd)
}

func Stat(cmd *cobra.Command, args []string) {
	res, err := stat(args[0])
	if err != nil {
		fmt.Printf("get file stat failed:\n\t %s\n", err.Error())
		log.Fatal("get file stat failed", log.Err(err))
	}

	fileType := "directory"
	if !res.IsDirectory {
		fileType = "file"
	}
	fmt.Printf("fileName: %s\n", res.FileName)
	fmt.Printf("fileSize(B): %d\n", res.Filesize)
	fmt.Printf("fileType: %s\n", fileType)
	fmt.Printf("createAt: %s\n", res.CreateTime)
	fmt.Printf("updateAt: %s\n", res.UpdateTime)


}

func stat(filepath string) (*pb1.StatResponse, error) {
	conn, err := grpc.Dial(config.Cfg.NamenodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	
	c := pb1.NewNameNodeClient(conn)
	req := &pb1.StatRequest{RemoteFilePath: filepath}
	res, err := c.Stat(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("get from namenode failed: %w", err)
	}

	return res, nil
}