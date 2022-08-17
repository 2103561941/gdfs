package server

import (
	"context"
	"fmt"

	"github.com/cyb0225/gdfs/internal/namenode/tree"
	"github.com/cyb0225/gdfs/pkg/log"
	pb1 "github.com/cyb0225/gdfs/proto/namenode"
	pb2 "github.com/cyb0225/gdfs/proto/datanode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *Server) Delete(ctx context.Context, req *pb1.DeleteRequest) (*pb1.DeleteResponse, error) {
	filepath := req.RemoteFilePath
	// delete file in file tree
	node, err := s.tree.Delete(filepath)
	if err != nil {
		log.Info("delete file failed", log.String("file", filepath), log.Err(err))
		return nil, fmt.Errorf("delete file failed: %w", err)
	}

	// get all node
	nodes := s.tree.GetChildrenNode(node)
	// get filekeys
	filekeys := getAllFileKeys(nodes)

	// delete filekey in cache
	// delete files in datanode
	for i := 0; i < len(filekeys); i++ {
		backups := s.cache.Get(filekeys[i]).Backups
		for j := 0; j < len(backups); j++ {
			sendDeleteMsg(backups[j], filekeys[i])
		}
		s.cache.Delete(filekeys[i])
	}

	log.Info("delete file success", log.String("file", filepath))
	res := &pb1.DeleteResponse{}
	return res, nil
}

// get all filekey.
func getAllFileKeys(nodes []*tree.Node) []string {
	// get all filekeys.
	var filekeys []string

	for i := 0; i < len(nodes); i++ {
		if !nodes[i].IsDirectory() {
			filekeys = append(filekeys, nodes[i].FileKeys...)
		}
	}

	return filekeys
}

func sendDeleteMsg(addr string, filekey string) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("connect to datanode failed", log.String("datanode", addr), log.Err(err))
		return
	}
	defer conn.Close()

	c := pb2.NewDataNodeClient(conn)
	req := &pb2.DeleteRequest{Filekey: filekey}

	if _, err = c.Delete(context.Background(), req); err != nil {
		log.Error("get from namenode failed", log.Err(err))
		return
	}
}
