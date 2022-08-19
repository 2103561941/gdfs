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

var mkdirCmd = &cobra.Command{
	Use:   "mkdir",
	Short: "make directory in remote file system",
	Args: func(cmd *cobra.Command, args []string) error {
		return equalNumArgs(1, args)
	},
	Run: Mkdir,
}

func init() {
	rootCmd.AddCommand(mkdirCmd)
}

func Mkdir(cmd *cobra.Command, args []string) {
	if _, err := mkdir(args[0]); err != nil {
		fmt.Printf("make directory failed:\n\t %s\n", err.Error())
		log.Fatal("make directory failed", log.Err(err))
	}

	fmt.Println("make directory success!")
	log.Info("make directroy success!")
}

func mkdir(filepath string) (*pb1.MkdirResponse, error) {
	conn, err := grpc.Dial(config.Cfg.NamenodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	c := pb1.NewNameNodeClient(conn)
	res, err := c.Mkdir(context.Background(), &pb1.MkdirRequset{RemoteFilePath: filepath})
	if err != nil {
		return nil, fmt.Errorf("get from namenode failed: %w", err)
	}
	return res, nil
}
