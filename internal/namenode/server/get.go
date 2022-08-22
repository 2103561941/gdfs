package server

import (
	"context"
	"fmt"

	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// Get addresses of datanodes which stored this file's chunks and backups.
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	filepath := req.RemoteFilePath
	node, err := s.tree.Get(filepath)
	if err != nil {
		return nil, fmt.Errorf("get file %s failed: %w", filepath, err)
	}

	// filekeys is the file chunks' key
	filekeys := node.FileKeys
	chunks := make([]*pb.Chunk, len(filekeys))
	// Get backups form filekeys.
	for i, filekey := range filekeys {
		chunk := &pb.Chunk{
			FileKey: filekey,
		}

		// Get addressed of datanodes which stored this file(filekey).
		backups := s.cache.Get(filekey)
		aliveBackups := []string{}
		log.Debugf("backups is %+v", backups)
		for j := 0; j < len(backups); j++ {
			if ok := s.alive.IsAlive(backups[j]); ok {
				aliveBackups = append(aliveBackups, backups[j])
			}
		}
		log.Debugf("aliveBackups is %+v", aliveBackups)

		chunk.Backups = aliveBackups
		chunks[i] = chunk
	}

	res := &pb.GetResponse{
		Chunks: chunks,
	}

	return res, nil
}
