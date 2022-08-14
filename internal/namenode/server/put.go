package server

import (
	"context"
	"fmt"
	"log"

	"github.com/cyb0225/gdfs/internal/namenode/tree"
	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// choose which datanode to store this file
// and create file chuncks' key, but not push it to the cache. It needs datanode report.

// Put the node into directory tree, but don't put it to cache,
// 	cache should be pushed when datanode already stored it.
// Prevent the problem that I have saved the file in the namenode,
// 	but there is no inconsistency of the uploaded data in the datanode.
// When user can get this file in directory tree, but cannot find it in cache,
// 	so he can not find which datanode stored it. It guaranteed file storage security.
// And I can also remove the lost file by starting asynchronous or multithreading.

// get datanode infomation
func (s *Server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	log.Println("into Put function")

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

	res := &pb.PutResponse{}

	return res, nil
}
