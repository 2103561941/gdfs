package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/cyb0225/gdfs/proto/datanode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *Server) Put(stream pb.DataNode_PutServer) error {

	// get filename, open/create file
	res, err := stream.Recv()
	if err != nil {
		return err
	}
	filekey := res.Filekey
	adds := res.Adds

	fd, err := os.OpenFile("./storage/tmp/"+filekey, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("cannot create file %s: %w", filekey, err)
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
			return fmt.Errorf("write to file %s failed: %w", filekey, err)
		}
	}

	// report to namenode.
	go func ()  {
		if err := s.report.FileReport(filekey); err != nil {
			log.Println("report filekey to namenode failed: %w", filekey)
		}		
	}()

	if len(adds) == 0 {
		return stream.SendAndClose(&pb.PutResponse{})
	}

	// asynchronous process
	// write to another datanodes.
	go func ()  {
		fd, err := os.Open(filekey)
		if err != nil {
			log.Printf("open file %s failed: %s", filekey, err.Error())
			return
		}
		defer fd.Close()
		r := bufio.NewReader(fd)

		for i := 0; i < len(adds); i++ {
			if err := putdata(filekey, r, adds[1:]); err != nil {
				log.Printf("put file %s to datanode failed: %s\n", filekey, err.Error())
				continue
			}
			break
		}
		
	}()


	return stream.SendAndClose(&pb.PutResponse{})
}

// put data to datanode
// add[0] stored the address which will be visited, and adds[1:] stored the other backups' address.
func putdata(filekey string, r io.Reader, adds []string) error {

	// if put one datanode failed, then try to put to next backups.
	// at the same time
	address := adds[0]

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("connect to server failed: %w", err)
	}
	defer conn.Close()

	c := pb.NewDataNodeClient(conn)
	stream, err := c.Put(context.Background())
	if err != nil {
		return fmt.Errorf("get stream failed: %w", err)
	}

	// send basic information to datanode.
	if err := stream.Send(&pb.PutRequest{Filekey: filekey, Adds: adds[1:]}); err != nil {
		return fmt.Errorf("send basic information to datanode %s failed: %w", address, err)
	}

	buf := make([]byte, 1024) //chunk size can divide it.   chunksize mod bufsize = 0
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			// size of buf is 0
			break
		}
		if err != nil {
			return fmt.Errorf("read filebytes from file %s failed: %w", filekey, err)
		}

		// fmt.Println("send buf: ", string(buf))
		if err := stream.Send(&pb.PutRequest{Databytes: buf[:n]}); err != nil {
			return fmt.Errorf("send basic information to datanode %s failed: %w", address, err)
		}
	}

	if _, err = stream.CloseAndRecv(); err != nil {
		return fmt.Errorf("close client stream failed: %w", err)
	}

	return nil
}
