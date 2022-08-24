package server

import (
	"context"
	"fmt"

	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// Rename
func (s *Server) Rename(ctx context.Context, req *pb.RenameRequest) (*pb.RenameResponse, error) {
	src := req.RenameSrcPath
	dest := req.RenameDestPath

	if err := s.tree.Rename(src, dest); err != nil {
		log.Info("rename file failed", log.String("src", src), log.String("dest", dest), log.Err(err))
		return nil, fmt.Errorf("rename file failed: %w", err)
	}

	if err := s.tree.Per.Rename(src, dest); err != nil {
		log.Error("write file tree log failed", log.String("method", "rename"), log.Err(err))
	}
	res := &pb.RenameResponse{}
	return res, nil 
}