package server

import (
	"context"

	pb "github.com/cyb0225/gdfs/proto/namenode"
)

func (s *Server) FileReport(ctx context.Context, req *pb.FileReportRequest) (*pb.FileReportResponse, error) {
	address := ctx.Value("address").(string)
	// log.Debug("filereport", log.String("datanode", address))
	filekey := req.Filekey

	res := &pb.FileReportResponse{
		IsExist: true,
	}
	// stored the mapping between filekey and address, then next time, use can use get to find the filekeys' datanode.
	if err := s.cache.Put(filekey, address); err != nil {
		res.IsExist = false
	}

	return res, nil
}
