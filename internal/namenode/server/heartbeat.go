// datanode heartbeat.

package server

import (
	"context"

	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// Datanode heartbeat.
// Make sure datanode is alive.
func (s *Server) HeartBeat(ctx context.Context, req *pb.HeartBeatRequset) (*pb.HeartBeatResponse, error) {
	address := ctx.Value("address").(string)
	s.alive.Update(address)
	res := &pb.HeartBeatResponse{}
	return res, nil
}
