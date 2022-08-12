// datanode heartbeat.

package server

import (
	"context"

	pb "github.com/cyb0225/gdfs/proto/namenode"
)

func (s *Server) HeartBeat(ctx context.Context, req *pb.HeartBeatRequset) (*pb.HeartBeatResponse, error) {
	address := req.Address
	s.alive.Update(address)
	res := &pb.HeartBeatResponse{
		Ack: 1,
	}
	return res, nil
}
