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
	// log.Debug("heartbeat", log.String("datanode", address), log.Int64("capacity", req.Cap))
	s.alive.Update(address, req.Cap)
	res := &pb.HeartBeatResponse{}
	return res, nil
}
