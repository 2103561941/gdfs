package server

import (
	"context"

	pb "github.com/cyb0225/gdfs/proto/namenode"
)

func (s *Server) FileReport(ctx context.Context, req *pb.FileReportRequest) (*pb.FileReportRequest, error) {

	address := req.Addr
	filekey := req.Filekey

	// stored the mapping between filekey and address, then next time, use can use get to find the filekeys' datanode.
	s.cache.Put(filekey, address)

	res := &pb.FileReportRequest{}
	return res, nil
}
