package report

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/cyb0225/gdfs/proto/namenode"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = viper.GetString("namenodeAddr")
)

type Report struct {}

func NewReport()*Report {
	return &Report{}
}

func (r *Report)HeartBeat() {
	fmt.Println("start heartbeat...")
	for {
		if err := heartbeat(); err != nil {
			// namenode may closed, in this way, datanode can choose another namenode.
			log.Printf("cannot connect to namenode: %s, please have a check: %s\n", addr, err.Error())
			os.Exit(1)
		}
		
		time.Sleep(time.Second * 20)
	}
}


func heartbeat() error {

	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("connect to namenode[%s] failed: %w", addr, err)
	}

	defer conn.Close()

	c := pb.NewNameNodeClient(conn)
	res, err := c.HeartBeat(context.Background(), &pb.HeartBeatRequset{Address: "127.0.0.1:50052"})
	if err != nil {
		return fmt.Errorf("get datanodes' information failed: %w", err)
	}

	fmt.Printf("%+v\n", res)

	return nil

}


// report file to namenode cache
func (r *Report)FileReport() {

}