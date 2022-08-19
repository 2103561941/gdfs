// datanode heartbeat.

package server

import (
	"context"
	"fmt"

	"github.com/cyb0225/gdfs/pkg/log"
	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// Datanode heartbeat.
// Make sure datanode is alive.
func (s *Server) HeartBeat(ctx context.Context, req *pb.HeartBeatRequset) (*pb.HeartBeatResponse, error) {
	address := ctx.Value("address").(string)
	if err := s.cache.Update(address); err != nil {
		log.Error("update address failed", log.String("datanode", address), log.Err(err))
		return nil, fmt.Errorf("update address failed: %w", err)
	}
	res := &pb.HeartBeatResponse{}
	return res, nil
}
