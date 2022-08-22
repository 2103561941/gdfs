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
	// Delete file in file tree
	node, err := s.tree.Delete(filepath)
	if err != nil {
		log.Info("delete file failed", log.String("file", filepath), log.Err(err))
		return nil, fmt.Errorf("delete file failed: %w", err)
	}

	// Get all node
	// Delete file may be a directory, so need to find and delete all its subfiles.
	nodes := s.tree.GetChildrenNode(node)
	filekeys := getAllFileKeys(nodes)

	// Delte filekey in cache. 
	// At the same time, get addresses of datanodes, and send delete msg to them. 
	for i := 0; i < len(filekeys); i++ {
		backups := s.cache.Get(filekeys[i])
		for j := 0; j < len(backups); j++ {
			// Check if the datanode is alived.
			if ok := s.alive.IsAlive(backups[i]); ok {
				sendDeleteMsg(backups[j], filekeys[i])
			}
		}
	}

	_ = s.tree.Per.Delete(filepath)
	log.Info("delete file success", log.String("file", filepath))
	res := &pb1.DeleteResponse{}
	return res, nil
}

// Get filekeys from a slice of file Node.
func getAllFileKeys(nodes []*tree.Node) []string {
	var filekeys []string

	for i := 0; i < len(nodes); i++ {
		if !nodes[i].IsDirectory() {
			filekeys = append(filekeys, nodes[i].FileKeys...)
		}
	}

	return filekeys
}

// Connect to datanode and Require it to delete relate file.
// But namenode don't care about if the datanode delete the file successfully.
// Namenode only needs to delete the file in file tree.
// In this way, the file is deleted to the outside world.
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
