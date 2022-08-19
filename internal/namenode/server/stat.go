package server

import (
	"context"
	"fmt"

	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// Return the stat of a file.
func (s *Server) Stat(ctx context.Context, req *pb.StatRequest) (*pb.StatResponse, error) {
	filepath := req.RemoteFilePath
	node, err := s.tree.Stat(filepath)
	if err != nil {
		log.Info("get file stat failed", log.String("file", filepath), log.Err(err))
		return nil, fmt.Errorf("get file stat failed: %w", err)
	}

	res := &pb.StatResponse{
		FileName:    node.FileName,
		Filesize:    node.FileSize,
		IsDirectory: node.IsDirectory(),
		UpdateTime:  node.UpdateTime.Format("2006-01-02 15:04:05.000"),
		CreateTime:  node.CreateTime.Format("2006-01-02 15:04:05.000"),
	}

	log.Info("get file stat success", log.String("file", filepath))
	return res, nil
}
