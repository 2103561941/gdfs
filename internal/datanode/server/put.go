package server

import (
	"bufio"
	"fmt"
	"io"
	"os"

	pb "github.com/cyb0225/gdfs/proto/datanode"
)

func (s *Server) Put(stream pb.DataNode_PutServer) error {

	// get filename, open/create file
	res, err := stream.Recv()
	if err != nil {
		return err
	}
	filename := res.Filename
	fd, err := os.OpenFile("./storage/tmp/"+filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("cannot create file %s: %w", filename, err)
	}
	defer fd.Close()
	w := bufio.NewWriter(fd)

	// write to
	for {
		res, err = stream.Recv()

		if err == io.EOF {
			w.Flush()
			break
		}
		if err != nil {
			return err
		}

		if _, err := w.Write(res.Databytes); err != nil {
			return fmt.Errorf("write to file %s failed: %w", filename, err)
		}
	}

	// TODO:
	// write to another datanodes.
	// report to namenode.

	return stream.SendAndClose(&pb.PutResponse{})
}
