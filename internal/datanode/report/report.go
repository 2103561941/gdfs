package report

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/cyb0225/gdfs/proto/namenode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	namenodeAddr = "127.0.0.1:50051"
	addr         = "127.0.0.1:50052"
)

type Report struct{}

func NewReport() *Report {
	return &Report{}
}

func (r *Report) HeartBeat() {
	fmt.Println("start heartbeat...")
	for {
		if err := heartbeat(); err != nil {
			// namenode may closed, in this way, datanode can choose another namenode.
			log.Printf("cannot connect to namenode: %s, please have a check: %s\n", namenodeAddr, err.Error())
			os.Exit(1)
		}

		time.Sleep(time.Second * 20)
	}
}

func heartbeat() error {
	conn, err := grpc.Dial(namenodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("connect to namenode[%s] failed: %w", namenodeAddr, err)
	}

	defer conn.Close()

	c := pb.NewNameNodeClient(conn)
	if _, err := c.HeartBeat(context.Background(), &pb.HeartBeatRequset{Addr: addr}); err != nil {
		return fmt.Errorf("get namenode' heartbeat failed: %w", err)
	}

	return nil
}

// report file to namenode cache
func (r *Report) FileReport(filekey string) error {
	conn, err := grpc.Dial(namenodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("connect to namenode[%s] failed: %w", namenodeAddr, err)
	}

	defer conn.Close()

	c := pb.NewNameNodeClient(conn)
	req := &pb.FileReportRequest{
		Filekey: filekey,
		Addr:    addr,
	}
	if _, err := c.FileReport(context.Background(), req); err != nil {
		return fmt.Errorf("get namenode' filereport failed: %w", err)
	}
	return nil
}
