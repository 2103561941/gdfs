package server

import (
	"context"

	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// When datanode stored a file successfully or restart, it will report the file to namenode. 
func (s *Server) FileReport(ctx context.Context, req *pb.FileReportRequest) (*pb.FileReportResponse, error) {
	address := ctx.Value("address").(string)
	filekey := req.Filekey
	s.cache.Put(filekey, address)

	res := &pb.FileReportResponse{}
	return res, nil
}
