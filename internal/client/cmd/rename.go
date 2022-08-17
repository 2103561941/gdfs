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

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "rename file in remote file system",
	Args: func(cmd *cobra.Command, args []string) error {
		return equalNumArgs(2, args)
	},
	Run: Rename,
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

func Rename(cmd *cobra.Command, args []string) {
	if _, err := rename(args[0], args[1]); err != nil {
		log.Fatal("rename failed", log.Err(err))
	}

	log.Info("rename success!")
}

func rename(src string, dest string) (*pb1.RenameResponse, error) {
	conn, err := grpc.Dial(config.Cfg.NamenodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	c := pb1.NewNameNodeClient(conn)
	req := &pb1.RenameRequest{
		RenameSrcPath:  src,
		RenameDestPath: dest,
	}
	res, err := c.Rename(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("get from namenode failed: %w", err)
	}
	return res, nil
}
