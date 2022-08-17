// datanode heartbeat.

package server

import (
	"context"

	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/namenode"
)

func (s *Server) HeartBeat(ctx context.Context, req *pb.HeartBeatRequset) (*pb.HeartBeatResponse, error) {
	address := ctx.Value("address").(string)
	// log.Debug("heartbeat", log.String("datanode", address))

	s.alive.Update(address)
	if ok := s.alive.IsAlive(address); !ok {
		log.Error("update address failed", log.String("datanode", address))
	}
	res := &pb.HeartBeatResponse{}
	return res, nil
}
