package server

import (
	"context"
	"os"

	"github.com/cyb0225/gdfs/internal/datanode/config"
	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/datanode"
)

func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	filekey := req.Filekey

	// only need to make sure the connect successfully,
	// and sync delete the  filekey, don't need to return whether the delection is successful.
	go delete(filekey)
	
	res := &pb.DeleteResponse{}
	return res, nil
}


func delete(filekey string) {
	// delete local file
	if err := os.Remove(config.Cfg.StoragePath + filekey); err != nil {
		log.Error("delete filekey failed", log.String("filekey", filekey), log.Err(err))
	}
}