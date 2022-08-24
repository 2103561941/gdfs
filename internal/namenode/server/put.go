package server

import (
	"context"
	"fmt"

	"github.com/cyb0225/gdfs/internal/namenode/tree"
	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// Choose which datanodes to store this file
// And create file chunks' key, but not put it to the cache, it depends on datanode reporting.
func (s *Server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	filepath := req.RemoteFilePath
	filesize := req.Filesize

	node := tree.NewNode(
		tree.IsDirectory(false),
		tree.SetFilePath(filepath),
		tree.SetFileSize(filesize),
	)

	// Create keys of file's partitions.
	if err := node.CreateFileKeys(s.chunkSize); err != nil {
		return nil, fmt.Errorf("failed to create file keys: %w", err)
	}

	if err := s.tree.Put(filepath, node); err != nil {
		return nil, err
	}

	chunks := make([]*pb.Chunk, len(node.FileKeys))
	// search datanode to store backups
	for i := 0; i < len(node.FileKeys); i++ {
		addrs, err := s.alive.LoadBalance(s.backups)
		// Put file error
		// There are not enough datanode to store this file.
		// Then delete node in file tree.
		if err != nil || len(addrs) == 0 {
			if _, err := s.tree.Delete(filepath); err != nil {
				log.Error("delete file failed", log.String("file", filepath), log.Err(err))
			}
			return nil, err
		}

		chunk := &pb.Chunk{
			Backups: addrs,
			FileKey: node.FileKeys[i],
		}
		chunks[i] = chunk
	}

	if err := s.tree.Per.Put(node); err != nil {
		log.Error("write file tree log failed", log.String("method", "put"), log.Err(err))
	}
	res := &pb.PutResponse{Chunks: chunks}
	return res, nil
}
