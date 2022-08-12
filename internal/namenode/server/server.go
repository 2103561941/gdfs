package server

import "gdfs/internal/datanode/server"

var _ server.Server = (*Server)(nil)

type Server struct {
	
}