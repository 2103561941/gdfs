package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/cyb0225/gdfs/internal/namenode/tree"
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


	if node.FileMeta.FileType == tree.Direcotry {
		return nil, fmt.Errorf("file: %s is a directory\n", filepath)	
	}

	// keys is the file chunks' key
	keys := node.FileMeta.FileKeys
	res := &pb.GetResponse{
		Chunks: make([]*pb.Chunk, len(keys)),
	}

	// get backups form keys
	for i, v := range keys {
		// get backups' datanode address
		backups := s.cache.Get(v)
		if backups == nil {
			return nil, errors.New("file is not exist")
		}

		// check is datanode alive
		for _, b := range backups.Backups {
			if ok := s.alive.IsAlive(b.Address); ok {
				if res.Chunks[i].Backups == nil {
					res.Chunks[i].Backups = make([]*pb.Backup, 1)
				}
				res.Chunks[i].Backups = append(res.Chunks[i].Backups, &pb.Backup{FileKey: b.Address})
			}	
		}
		// it turns out that there is no datanode store this file chunk. file is lost.
		if res.Chunks[i].Backups == nil {
			return nil, errors.New("file is lost")
		}
	}

	return res, nil
}