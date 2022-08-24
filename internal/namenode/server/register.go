package server

import (
	"context"

	pb "github.com/cyb0225/gdfs/proto/namenode"
)

// Datanode Register.
// Datanode register is to check if the datanode is still in cache.
// If datanode is still in cache and don't delete, then delete the old cache and put new filekeys to the cache.
func (s *Server) Register(ctx context.Context, req *pb.RegisterRequset) (*pb.RegisterResponse, error) {
	//TODO:
	// address := ctx.Value("address").(string)
	// s.alive.Update(address)
	res := &pb.RegisterResponse{}
	return res, nil
}