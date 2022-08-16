package server

import (
	"net"

	"github.com/cyb0225/gdfs/internal/datanode/config"
	"github.com/cyb0225/gdfs/internal/datanode/report"
	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/datanode"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedDataNodeServer
	report *report.Report
}

func NewServer() *Server {
	return &Server{
		report: report.NewReport(),
	}
}

// start rpc server
func RunServer() error {
	lis, err := net.Listen("tcp", ":"+config.Cfg.Addr.Port)
	if err != nil {
		return err
	}
	datanodeServer := NewServer()
	s := grpc.NewServer()
	pb.RegisterDataNodeServer(s, datanodeServer)

	log.Info("server start listening", log.String("port", config.Cfg.Addr.Port))

	go func() {
		if err = s.Serve(lis); err != nil {
			log.Fatalf("start server failed: %s", err.Error())
		}
	}()

	// restart file report
	if err := datanodeServer.report.RestartFileReport(); err != nil {
		return err
	}

	datanodeServer.report.HeartBeat()


	return nil
}
