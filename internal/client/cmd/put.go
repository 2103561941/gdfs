package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb2 "github.com/cyb0225/gdfs/proto/datanode"
	pb1 "github.com/cyb0225/gdfs/proto/namenode"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "put local file to remote file system",
	Args: func(cmd *cobra.Command, args []string) error {
		return equalNumArgs(2, args)
	},
	Run: Put,
}

func init() {
	rootCmd.AddCommand(putCmd)
}

func Put(cmd *cobra.Command, args []string) {
	localFilePath := args[0]
	remoteFilePath := args[1] // choose a place to store this file

	fd, err := os.Open(localFilePath)
	if err != nil {
		log.Fatalf("failed to open file: %s: %s\n", localFilePath, err.Error())
	}
	defer fd.Close()

	fileinfo, err := fd.Stat()
	if err != nil {
		log.Fatalf("failed to get file: %s stat: %s\n", localFilePath, err.Error())
	}
	// get bytes, should transform to KB
	filesize := (fileinfo.Size())
	fmt.Println("filesize: ", filesize)

	res, err := put(remoteFilePath, filesize)
	if err != nil {
		log.Fatalf("get datanode information from namenode failed: %s\n", err.Error())
	}

	r := bufio.NewReader(fd)

	for i := 0; i < len(res.Chunks); i++ {
		filename := res.Chunks[i].FileKey
		backups := res.Chunks[i].Backups
		if err := putdata(filename, r, backups); err != nil {
			log.Fatalf("put file %s to datanode failed: %s\n", filename, err.Error())
		}
	}

	fmt.Printf("client get server: %+v", res)
}

// get datanode information from namenode
func put(filepath string, filesize int64) (*pb1.PutResponse, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to server failed: %w", err)
	}
	defer conn.Close()

	c := pb1.NewNameNodeClient(conn)
	req := &pb1.PutRequest{
		RemoteFilePath: filepath,
		Filesize:       filesize,
	}
	res, err := c.Put(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("get from %s server failed: %w", addr, err)
	}

	return res, nil
}

// put data to datanode
func putdata(filename string, r io.Reader, adds []string) error {

	// if put one datanode failed, then try to put to next backups.
	// at the same time
	hasError := true
	for i := 0; i < len(adds); i++ {
		address := adds[i]

		conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			// return fmt.Errorf("connect to server failed: %w", err)
			log.Printf("connect to server failed: %s\n", err.Error())
			conn.Close()
			continue
		}
		
		c := pb2.NewDataNodeClient(conn)
		stream, err := c.Put(context.Background())
		if err != nil {
			// return fmt.Errorf("get stream failed: %w", err)
			log.Printf("get stream failed: %s\n", err.Error())
			conn.Close()
			continue
		}
			
		// send basic information to datanode.
		if err := stream.Send(&pb2.PutRequest{Filename: filename, Adds: adds[i + 1:]}); err != nil {
			log.Printf("send basic information to datanode %s failed: %s\n", address, err.Error())
			conn.Close()
			continue
		} 
		buf := make([]byte, 1024)
		sum := 0 // stored the bytes that read. 
		for {
			n, err := r.Read(buf)
			if err == io.EOF {
				// size of buf is 0
				break
			}
			if err != nil {
				log.Printf("read filebytes from file %s failed: %s\n", filename, err.Error() )
				conn.Close()
				continue
			}

			sum += n
			if sum >= 1024 * 1024 { // every chunk's size
				break
			}

			if err := stream.Send(&pb2.PutRequest{Databytes: buf[:n]}); err != nil {
				log.Printf("send basic information to datanode %s failed: %s\n", address, err.Error())
				conn.Close()
				continue
				
			}

		}
			
		conn.Close()
		hasError = false // it truns out that put have successd at least once
		break
	}
	if hasError {
		return fmt.Errorf("cannot put file to any datanode")
	}

	return nil
}
