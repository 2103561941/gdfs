package server

import (
	"bufio"
	"fmt"
	"io"
	"os"

	pb "github.com/cyb0225/gdfs/proto/datanode"
)

func (s *Server) Get(req *pb.GetRequset, stream pb.DataNode_GetServer) error {
	filename := req.Filename
	fd, err := os.Open("./storage/tmp")
	if err != nil {
		return fmt.Errorf("cannot open file %s : %w", filename, err)
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
			return fmt.Errorf("read file %s failed: %w", filename, err)
		}

		if err := stream.Send(&pb.GetResponse{
			Databytes: buf[:n],
		}); err != nil {
			return fmt.Errorf("send databytes failed: %w", err)
		}
	}
	return nil
}
