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

	if err := s.tree.Put(filepath, &tree.Node{
		FileMeta: &tree.Meta{
			FileType: tree.NormalFile,
			Path: filepath,
		},
		Children: make([]*tree.Node, 0),
	}); err != nil {
		return nil, err
	}

	return &pb.PutResponse{}, nil
}
