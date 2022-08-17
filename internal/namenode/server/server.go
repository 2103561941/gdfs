package server

import (
	"fmt"
	"net"

	"github.com/cyb0225/gdfs/internal/namenode/alive"
	"github.com/cyb0225/gdfs/internal/namenode/cache"
	"github.com/cyb0225/gdfs/internal/namenode/config"
	"github.com/cyb0225/gdfs/internal/namenode/tree"
	"github.com/cyb0225/gdfs/internal/pkg/middleware"
	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/namenode"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
)

var _ pb.NameNodeServer = (*Server)(nil)

type Server struct {
	pb.UnimplementedNameNodeServer
	tree  *tree.Tree
	cache *cache.Cache
	alive *alive.Alive
}

func newServer() *Server {
	return &Server{
		tree:  tree.NewTree(),
		cache: cache.NewCache(),
		alive: alive.NewAlive(),
	}
}

// start rpc server
func RunServer() error {
	lis, err := net.Listen("tcp", ":"+config.Cfg.Addr.Port)
	if err != nil {
		return err
	}

	logger := log.ZapLogger()
	if logger == nil {
		return fmt.Errorf("log not init, can not get zap logger")
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.UneryRecovery(),
				grpc_ctxtags.UnaryServerInterceptor(),
				middleware.UnaryServerInterceptor([]string{"HeartBeat"}),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				middleware.StreamRecovery(),
				grpc_ctxtags.StreamServerInterceptor(),
				middleware.StreamServerInterceptor(nil),
			)),
	)

	pb.RegisterNameNodeServer(s, newServer())

	log.Info("server start listening", log.String("port", config.Cfg.Addr.Port))
	if err = s.Serve(lis); err != nil {
		return err
	}

	return nil
}
