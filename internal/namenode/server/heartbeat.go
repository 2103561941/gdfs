// datanode heartbeat.

package server

import (
	"context"

	pb "github.com/cyb0225/gdfs/proto/namenode"
)

func (s *Server) HeartBeat(ctx context.Context, req *pb.HeartBeatRequset) (*pb.HeartBeatResponse, error) {
	address := req.Addr

	s.alive.Update(address)
	s.alive.IsAlive(address)
	res := &pb.HeartBeatResponse{}
	return res, nil
}
