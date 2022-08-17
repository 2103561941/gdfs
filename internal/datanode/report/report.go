package report

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cyb0225/gdfs/internal/datanode/config"
	"github.com/cyb0225/gdfs/internal/pkg/middleware"
	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/namenode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Report struct{}

func NewReport() *Report {
	return &Report{}
}

func (r *Report) HeartBeat() {
	log.Info("start heartbeat")
	for {
		if err := heartbeat(); err != nil {
			// namenode may closed, in this way, datanode can choose another namenode.
			log.Fatal("cannot connect to namenode", log.String("namenode", config.Cfg.NamenodeAddr), log.Err(err))
			os.Exit(1)
		}

		time.Sleep(time.Second * time.Duration(config.Cfg.HeartBeatTime))
	}
}

func heartbeat() error {
	addr := config.Cfg.Addr.IP + ":" + config.Cfg.Addr.Port

	conn, err := grpc.Dial(config.Cfg.NamenodeAddr, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(middleware.UnaryClientInterceptor(addr)))
	if err != nil {
		return err
	}

	defer conn.Close()
	
	c := pb.NewNameNodeClient(conn)
	if _, err := c.HeartBeat(context.Background(), &pb.HeartBeatRequset{}); err != nil {
		return fmt.Errorf("get namenode' heartbeat failed: %w", err)
	}

	return nil
}

// report file to namenode cache
func (r *Report) FileReport(filekey string) error {
	addr := config.Cfg.Addr.IP + ":" + config.Cfg.Addr.Port

	conn, err := grpc.Dial(config.Cfg.NamenodeAddr, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(middleware.UnaryClientInterceptor(addr)),
	)
	if err != nil {
		return fmt.Errorf("connect to namenode[%s] failed: %w", config.Cfg.NamenodeAddr, err)
	}

	defer conn.Close()

	c := pb.NewNameNodeClient(conn)
	req := &pb.FileReportRequest{
		Filekey: filekey,
	}
	if _, err := c.FileReport(context.Background(), req); err != nil {
		return fmt.Errorf("get namenode' filereport failed: %w", err)
	}
	return nil
}

// Restart file reporting
// tell namenode, datanode stored files.
func (r *Report) RestartFileReport() error {
	fileInfos, err := os.ReadDir(config.Cfg.StoragePath)
	if err != nil {
		return fmt.Errorf("open directory failed: %w", err)
	}

	for _, file := range fileInfos {
		if err := r.FileReport(file.Name()); err != nil {
 			return err
		}
	}

	log.Info("restart file report success")

	return nil
}
