package server

import (
	"context"
	"fmt"

	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// Make directory.
func (s *Server) Mkdir(ctx context.Context, req *pb.MkdirRequset) (*pb.MkdirResponse, error) {

	filepath := req.RemoteFilePath
	node, err := s.tree.Mkdir(filepath);
	if err != nil {
		return nil, fmt.Errorf("make directory failed: %w", err)
	}

	_ = s.tree.Per.Mkdir(node)
	res := &pb.MkdirResponse{}
	return res, nil
}
