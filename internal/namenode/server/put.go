package server

import (
	"context"
	"fmt"

	"github.com/cyb0225/gdfs/internal/namenode/tree"
	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// choose which datanode to store this file
// and create file chunks' key, but not push it to the cache. It needs datanode report.

// Put the node into directory tree, but don't put it to cache,
// 	cache should be pushed when datanode already stored it.
// Prevent the problem that I have saved the file in the namenode,
// 	but there is no inconsistency of the uploaded data in the datanode.
// When user can get this file in directory tree, but cannot find it in cache,
// 	so he can not find which datanode stored it. It guaranteed file storage security.
// And I can also remove the lost file by starting asynchronous or multithreading.

// get datanode infomation
func (s *Server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {

	filepath := req.RemoteFilePath
	filesize := req.Filesize

	node := tree.NewNode(
		tree.IsDirectory(false),
		tree.SetFilePath(filepath),
		tree.SetFileSize(filesize),
	)

	if err := node.CreateFileKeys(); err != nil {
		return nil, fmt.Errorf("failed to create file keys: %w", err)
	}

	// error
	if err := s.tree.Put(filepath, node); err != nil {
		return nil, err
	}

	chunks := make([]*pb.Chunk, len(node.FileKeys))
	// search datanode to store backups
	for i := 0; i < len(node.FileKeys); i++ {
		adds, err := s.alive.Backup()
		// put file error
		if err != nil {
			if _, err := s.tree.Delete(filepath); err != nil {
				log.Error("delete file failed", log.String("file", filepath), log.Err(err))
			}
			return nil, err
		}

		log.Debugf("put adds %v", adds)

		chunk := &pb.Chunk{
			Backups: adds,
			FileKey: node.FileKeys[i],
		}
		chunks[i] = chunk
	}

	// create filekey in cache
	for i := 0; i < len(chunks); i++ {
		if err := s.cache.Create(chunks[i].FileKey); err != nil {
			log.Error("create filekey to cache failed", log.String("filekey", chunks[i].FileKey), log.Err(err))
			// delete filekeys in cache
			for j := 0; j < i; j++ {
				s.cache.Delete(chunks[j].FileKey)
			}
			return nil, fmt.Errorf("create filekey to cache failed") 
		}
	}

	res := &pb.PutResponse{Chunks: chunks}
	log.Debug("put datanodes' address success")

	return res, nil
}
