package server

import (
	"fmt"
	"net"

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
	report *report.Report // File report and heartbeat to namenode.

	storagePath string
}

type ServerConfig struct {
	IP          string
	Port        string
	StoragePath string
	NamenodeAddr  string
    HeartBeatTime int
}

func NewServer(cfg *ServerConfig) *Server {
	report := &report.Report{
		IP: cfg.IP,
		Port: cfg.Port,
		Addr: cfg.IP + ":" + cfg.Port,
		StoragePath: cfg.StoragePath,
		NamenodeAddr: cfg.NamenodeAddr,
		HeartBeatTime: cfg.HeartBeatTime,
	}

	return &Server{
		report:      report,
		storagePath: cfg.StoragePath,
	}
}

// start rpc server
func RunServer(cfg *ServerConfig) error {
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return err
	}
	datanodeServer := NewServer(cfg)

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
	log.Info("server start listening", log.String("port", cfg.Port))

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
