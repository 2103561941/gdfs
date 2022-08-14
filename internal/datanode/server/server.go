package server

import (
	"log"
	"net"

	"github.com/cyb0225/gdfs/internal/datanode/report"
	pb "github.com/cyb0225/gdfs/proto/datanode"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedDataNodeServer
	report *report.Report
}

func newServer() *Server {
	return &Server{
		report: report.NewReport(),
	}
}

// start rpc server
func RunServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	datanodeServer := newServer()
	s := grpc.NewServer()
	pb.RegisterDataNodeServer(s, datanodeServer)

	log.Printf("server start listening at %s", port)

	go func ()  {
		if err = s.Serve(lis); err != nil {
			log.Fatalf("start server failed: %s", err.Error())
		}
	}()	

	datanodeServer.report.HeartBeat()

	return nil
}
