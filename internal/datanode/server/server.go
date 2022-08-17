package server

import (
	"fmt"
	"net"

	"github.com/cyb0225/gdfs/internal/datanode/config"
	"github.com/cyb0225/gdfs/internal/datanode/report"
	"github.com/cyb0225/gdfs/internal/pkg/middleware"
	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/datanode"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
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

	logger := log.ZapLogger()
	if logger == nil {
		return fmt.Errorf("log not init, can not get zap logger")
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middleware.UneryRecovery(),
			grpc_ctxtags.UnaryServerInterceptor(),
			middleware.UnaryServerInterceptor(nil),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			middleware.StreamRecovery(),
			grpc_ctxtags.StreamServerInterceptor(),
			middleware.StreamServerInterceptor(nil),
		)),
	)
	
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
