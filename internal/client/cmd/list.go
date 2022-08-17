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
		log.Fatal("list dir failed", log.Err(err))
	}

	// show results
	log.Info("list dir success")
	for i := 0; i < len(res.FileInfos); i++ {
		log.Info("file stat",
			log.String("filename", res.FileInfos[i].FileName),
			log.Int64("filesize", res.FileInfos[i].Filesize),
			log.Bool("isDirectory", res.FileInfos[i].IsDirectory),
			log.String("createTime", res.FileInfos[i].CreateTime),
			log.String("updateTime", res.FileInfos[i].UpdateTime),
		)
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
