package server

import (
	"context"
	"fmt"

	pb "github.com/cyb0225/gdfs/proto/namenode"
)


func (s *Server) Mkdir(ctx context.Context, req *pb.MkdirRequset) (*pb.MkdirResponse, error) {

	filepath := req.RemoteFilePath
	if err := s.tree.Mkdir(filepath); err != nil {
		return nil, fmt.Errorf("make directory failed: %w", err)
	}

	res := &pb.MkdirResponse{}
	return res, nil 
}
