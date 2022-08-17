package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cyb0225/gdfs/internal/datanode/config"
	"github.com/cyb0225/gdfs/internal/pkg/middleware"
	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/datanode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *Server) Put(stream pb.DataNode_PutServer) error {

	// it promises that file descriptor must closed more early than groutine putdata.
	// because, put data needs to open file, so the file descriptor should be closed before groutine opened it.
	// and defer is a stack, so fd.Close() is execute before "ch <- 1".Then putdata groutine.
	ch := make(chan int)
	defer func ()  {
		ch <- 1	
	}()

	// get filename, open/create file
	res, err := stream.Recv()
	if err != nil {
		return err
	}
	filekey := res.Filekey
	adds := res.Adds

	fd, err := os.OpenFile(config.Cfg.StoragePath+filekey, os.O_CREATE|os.O_WRONLY, 0666)
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

	go func() {
		if err := s.report.FileReport(filekey); err != nil {
			log.Error("report filekey to namenode failed", log.String("filekey", filekey), log.Err(err))
		}
	}()

	// sync put data to other datanodes.
	if len(adds) == 0 {
		return stream.SendAndClose(&pb.PutResponse{})
	}

	
	// report to namenode.
	go func() {
		<- ch
		 
		fd, err := os.Open(config.Cfg.StoragePath + filekey)
		if err != nil {
			log.Error("open file failed", log.String("filekey", filekey), log.Err(err))
			return
		}
		defer fd.Close()
		r := bufio.NewReader(fd)

		for i := 0; i < len(adds); i++ {
			if err := putdata(filekey, r, adds[i:]); err != nil {
				log.Error("put file to datanode failed", log.String("datanode", adds[i]), log.String("filekey", filekey), log.Err(err))
				continue
			}
			break
		}
	}()

	if err := stream.SendAndClose(&pb.PutResponse{}); err != nil {
		return err
	}

	return nil
}

// put data to datanode
// add[0] stored the address which will be visited, and adds[1:] stored the other backups' address.
func putdata(filekey string, r io.Reader, adds []string) error {

	// if put one datanode failed, then try to put to next backups.
	// at the same time
	address := adds[0] // other datanode's ip 
	addr := config.Cfg.Addr.IP + ":" + config.Cfg.Addr.Port // this datanode's ip

	conn, err := grpc.Dial(address, 
						grpc.WithTransportCredentials(insecure.NewCredentials()),
						grpc.WithChainUnaryInterceptor(middleware.UnaryClientInterceptor(addr)))

	if err != nil {
		return fmt.Errorf("connect to server failed: %w", err)
	}
	defer conn.Close()

	log.Debug("put data to datanode", log.String("ip", addr), log.String("datanode", address))
	c := pb.NewDataNodeClient(conn)
	stream, err := c.Put(context.Background())
	if err != nil {
		return fmt.Errorf("get stream failed: %w", err)
	}

	// send basic information to datanode.
	if err := stream.Send(&pb.PutRequest{Filekey: filekey, Adds: adds[1:]}); err != nil {
		return fmt.Errorf("send basic information to datanode failed: %w", err)
	}

	buf := make([]byte, 1024) //chunk size can divide it.   chunksize mod bufsize = 0
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			// size of buf is 0
			break
		}
		if err != nil {
			return fmt.Errorf("read file failed: %w", err)
		}

		if err := stream.Send(&pb.PutRequest{Databytes: buf[:n]}); err != nil {
			return fmt.Errorf("send filestat to datanode failed: %w", err)
		}
	}

	if _, err = stream.CloseAndRecv(); err != nil {
		return fmt.Errorf("close client stream failed: %w", err)
	}

	return nil
}
