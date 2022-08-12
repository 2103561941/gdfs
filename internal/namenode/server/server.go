package server

import (
	"log"
	"net"

	"github.com/cyb0225/gdfs/internal/namenode/alive"
	"github.com/cyb0225/gdfs/internal/namenode/cache"
	"github.com/cyb0225/gdfs/internal/namenode/tree"
	pb "github.com/cyb0225/gdfs/proto/namenode"
	"google.golang.org/grpc"
)

var _ pb.NameNodeServer = (*Server)(nil)

type Server struct {
	pb.UnimplementedNameNodeServer
	tree *tree.Tree
	cache *cache.Cache
	alive *alive.Alive
}

func newServer() *Server {
	return &Server{
		tree: tree.NewTree(),
		cache: cache.NewCache(),
		alive: alive.NewAlive(),
	}
}

// start rpc server
func RunServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterNameNodeServer(s, newServer())

	log.Printf("server start listening at %s", port)
	if err = s.Serve(lis); err != nil {
		return err
	}

	return nil
}
