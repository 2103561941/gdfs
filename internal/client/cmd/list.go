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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list files in remote file system",
	Args: func(cmd *cobra.Command, args []string) error {
		return equalNumArgs(1, args)
	},
	Run: List,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func List(cmd *cobra.Command, args []string) {
	res, err := list(args[0])
	if err != nil {
		fmt.Printf("list dir failed:\n\t %s\n", err.Error())
		log.Fatal("list dir failed", log.Err(err))
	}

	// show results
	log.Info("list dir success!")
	if len(res.FileInfos) == 0 {
		fmt.Println("empty.")
		return
	}
	fmt.Println("name size(B) type createAt updateAt")

	for i := 0; i < len(res.FileInfos); i++ {
		file := res.FileInfos[i]
		fileType := "directory"
		if !file.IsDirectory {
			fileType = "file"
		}
		fmt.Printf("%s %d %s %s %s\n", file.FileName, file.Filesize, fileType, file.CreateTime, file.UpdateTime)
	}
}

func list(dirpath string) (*pb1.ListResponse, error) {
	conn, err := grpc.Dial(config.Cfg.NamenodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := pb1.NewNameNodeClient(conn)
	req := &pb1.ListRequest{RemoteDirPath: dirpath}
	res, err := c.List(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("get from namenode failed: %w", err)
	}

	return res, nil
}
