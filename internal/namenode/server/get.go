package server

import (
	"context"
	"fmt"

	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// get datanode infomation
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {

	filepath := req.RemoteFilePath
	node := s.tree.Get(filepath)
	if node == nil {
		return nil, fmt.Errorf("file: %s is not exist", filepath)
	}

	if node.IsDirectory() {
		return nil, fmt.Errorf("file: %s is a directory\n", filepath)
	}

	// keys is the file chunks' key
	keys := node.FileKeys

	chunks := make([]*pb.Chunk, len(keys))

	// get backups form keys
	for i, key := range keys {
		chunk := &pb.Chunk{
			FileKey: key,
		}
		chunk.FileKey = key

		chunk.Backups = make([]string, 0)
		// get backups' datanode address
		backups := s.cache.Get(key)
		if backups == nil || len(backups.Backups) == 0 {
			return nil, fmt.Errorf("file is not exist")
		}

		// check is datanode alive
		for _, address := range backups.Backups {
			if ok := s.alive.IsAlive(address); ok {
				chunk.Backups = append(chunk.Backups, address)
			}
		}
		// it turns out that there is no datanode store this file chunk. file is lost.
		if chunk.Backups == nil || len(chunk.Backups) == 0 {
			return nil, fmt.Errorf("file is lost")
		}

		chunks[i] = chunk
	}

	res := &pb.GetResponse{
		Chunks: chunks,
	}

	return res, nil
}
