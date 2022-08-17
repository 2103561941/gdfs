package server

import (
	"context"
	"fmt"

	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/namenode"
)

func (s *Server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	dirpath := req.RemoteDirPath
	nodes, err := s.tree.List(dirpath)
	if err != nil {
		log.Info("list dir failed", log.String("dirpath", dirpath), log.Err(err))
		return nil, fmt.Errorf("list dir failed: %w", err)
	}

	fileInfos := make([]*pb.StatResponse, len(nodes))
	for i := 0; i < len(nodes); i++ {
		stat := &pb.StatResponse{
			FileName:    nodes[i].FileName,
			Filesize:    nodes[i].FileSize,
			IsDirectory: nodes[i].IsDirectory(),
			UpdateTime:  nodes[i].UpdateTime.Format("2006-01-02 15:04:05.000"),
			CreateTime:  nodes[i].CreateTime.Format("2006-01-02 15:04:05.000"),
		}
		fileInfos[i] = stat
	}
	res := &pb.ListResponse{
		FileInfos: fileInfos,
	}

	log.Info("list dir success", log.String("dirpath", dirpath))
	return res, nil
}
