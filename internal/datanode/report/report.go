package report

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cyb0225/gdfs/internal/pkg/middleware"
	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/namenode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Report struct {
	NamenodeAddr  string
	HeartBeatTime int
	IP            string
	Port          string
	Addr          string
	StoragePath   string
}

func (r *Report) HeartBeat() {
	log.Info("start heartbeat")
	for {
		if err := r.heartbeat(); err != nil {
			// namenode may closed, in this way, datanode can choose another namenode.
			log.Fatal("cannot connect to namenode", log.String("namenode", r.NamenodeAddr), log.Err(err))
			os.Exit(1)
		}

		time.Sleep(time.Second * time.Duration(r.HeartBeatTime))
	}
}

func (r *Report) heartbeat() error {
	conn, err := grpc.Dial(r.NamenodeAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(middleware.UnaryClientInterceptor(r.Addr)))
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
	conn, err := grpc.Dial(r.NamenodeAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(middleware.UnaryClientInterceptor(r.Addr)),
	)
	if err != nil {
		return fmt.Errorf("connect to namenode[%s] failed: %w", r.NamenodeAddr, err)
	}

	defer conn.Close()

	c := pb.NewNameNodeClient(conn)
	req := &pb.FileReportRequest{
		Filekey: filekey,
	}

	_, err = c.FileReport(context.Background(), req)
	if err != nil {
		return fmt.Errorf("get namenode' filereport failed: %w", err)
	}

	// filekey not exist in namenode. delete it.

	return nil
}

// Restart file reporting
// tell namenode, datanode stored files.
func (r *Report) RestartFileReport() error {
	fileInfos, err := os.ReadDir(r.StoragePath)
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
