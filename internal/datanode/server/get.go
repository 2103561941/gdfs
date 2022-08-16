package server

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/cyb0225/gdfs/internal/datanode/config"
	pb "github.com/cyb0225/gdfs/proto/datanode"
)

func (s *Server) Get(req *pb.GetRequset, stream pb.DataNode_GetServer) error {
	filekey := req.Filekey
	fd, err := os.Open(config.Cfg.StoragePath + filekey)
	if err != nil {
		return fmt.Errorf("cannot open file %s : %w", filekey, err)
	}
	defer fd.Close()
	r := bufio.NewReader(fd)

	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		
		if err != nil {
			return fmt.Errorf("read file %s failed: %w", filekey, err)
		}

		if err := stream.Send(&pb.GetResponse{
			Databytes: buf[:n],
		}); err != nil {
			return fmt.Errorf("send databytes failed: %w", err)
		}
	}

	return nil
}
