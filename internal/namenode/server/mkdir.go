package server

import (
	"context"
	"fmt"

	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// Make directory.
func (s *Server) Mkdir(ctx context.Context, req *pb.MkdirRequset) (*pb.MkdirResponse, error) {

	filepath := req.RemoteFilePath
	node, err := s.tree.Mkdir(filepath)
	if err != nil {
		return nil, fmt.Errorf("make directory failed: %w", err)
	}

	if err := s.tree.Per.Mkdir(node); err != nil {
		log.Error("write file tree log failed", log.String("method", "mkdir"), log.Err(err))
	}
	res := &pb.MkdirResponse{}
	return res, nil
}
