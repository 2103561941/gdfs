package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// get datanode infomation
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Println("into Get function")

	filepath := req.RemoteFilePath
	node := s.tree.Get(filepath)
	if node == nil {
		return nil, fmt.Errorf("file: %s not exist", filepath)
	}

	if node.IsDirectory() {
		return nil, fmt.Errorf("file: %s is a directory\n", filepath)
	}

	// keys is the file chunks' key
	keys := node.FileKeys
	chunks := make([]*pb.Chunk, 0)

	// get backups form keys
	for i, key := range keys {
		// get backups' datanode address
		backups := s.cache.Get(key)
		if backups == nil || len(backups.Backups) == 0 {
			return nil, errors.New("file is not exist")
		}

		if chunks[i].Backups == nil {
			chunks[i].Backups = make([]*pb.Backup, 0)
		}
		chunks[i].FileKey = key

		// check is datanode alive
		for _, b := range backups.Backups {
			if ok := s.alive.IsAlive(b.Address); ok {
				chunks[i].Backups = append(chunks[i].Backups, &pb.Backup{Address: b.Address})
			}
		}
		// it turns out that there is no datanode store this file chunck. file is lost.
		if chunks[i].Backups == nil || len(chunks[i].Backups) == 0 {
			return nil, errors.New("file is lost")
		}
	}

	res := &pb.GetResponse{
		Chunks: chunks,
	}

	return res, nil
}
