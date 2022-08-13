package server

import (
	"context"
	"log"

	"github.com/cyb0225/gdfs/internal/namenode/tree"
	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// get datanode infomation
func (s *Server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	log.Println("into Put function")

	filepath := req.RemoteFilePath

	node := tree.NewNode(
		tree.IsDirectory(false),
		tree.SetFilePath(filepath),
	)
	
	if err := s.tree.Put(filepath, node); err != nil {
		return nil, err
	}

	return &pb.PutResponse{}, nil
}
