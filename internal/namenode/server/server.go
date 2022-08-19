package server

import (
	"fmt"
	"net"

	"github.com/cyb0225/gdfs/internal/namenode/cache"
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

	// file tree.
	// Stored relation of files.
	tree *tree.Tree

	// filekey and datanode cache.
	// Namnode does not record what files stored in datanode.
	// It relies on datanode file reporting to stored the relation between filekey and datanode.
	cache *cache.Cache

	backups int // defautl file backups
	chunkSize int64 // file size per file block
}

type ServerConfig struct {
	Port string
	Backups int 
	ChunkSize int64
	StoragePath string 
	Expired int
}

func newServer(cfg *ServerConfig) (*Server, error) {
	tree, err := tree.NewTree(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("new tree failed: %w", err)
	}

	cache := cache.NewCache(cfg.Expired)
	server := &Server{
		tree:  tree,
		cache: cache,
		chunkSize: cfg.ChunkSize,
		backups: cfg.Backups,
	}

	return server, nil
}

// Start rpc server.
// Add some middlewares. Such as recovery, interceptor log.
func RunServer(cfg *ServerConfig) error {
	lis, err := net.Listen("tcp", ":"+ cfg.Port)
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
				// middleware.UneryRecovery(),
				grpc_ctxtags.UnaryServerInterceptor(),
				middleware.UnaryServerInterceptor([]string{"HeartBeat"}),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				// middleware.StreamRecovery(),
				grpc_ctxtags.StreamServerInterceptor(),
				middleware.StreamServerInterceptor(nil),
			)),
	)

	server, err := newServer(cfg)
	if err != nil {
		return fmt.Errorf("new server failed: %w", err)
	}

	pb.RegisterNameNodeServer(s, server)

	log.Info("server start listening", log.String("port", cfg.Port))
	if err = s.Serve(lis); err != nil {
		return err
	}

	return nil
}
