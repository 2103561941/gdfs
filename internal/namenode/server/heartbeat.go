// datanode heartbeat.

package server

import (
	"context"
	"fmt"

	pb "github.com/cyb0225/gdfs/proto/namenode"
)

func (s *Server) HeartBeat(ctx context.Context, req *pb.HeartBeatRequset) (*pb.HeartBeatResponse, error) {
	address := req.Addr
	fmt.Println(address)

	s.alive.Update(address)
	s.alive.IsAlive(address)
	res := &pb.HeartBeatResponse{}
	return res, nil
}
