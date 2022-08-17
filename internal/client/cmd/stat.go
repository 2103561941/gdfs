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
		log.Fatal("get file stat failed", log.Err(err))
	}

	log.Info("get file stat success", 
			log.String("filename", res.FileName),
			log.Int64("filesize", res.Filesize),
			log.Bool("isDirectory", res.IsDirectory),
			log.String("createTime", res.CreateTime),
			log.String("updateTime", res.UpdateTime),
		)
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