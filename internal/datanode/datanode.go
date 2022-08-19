package datanode

import (
	"log"

	"github.com/cyb0225/gdfs/internal/datanode/config"
	"github.com/cyb0225/gdfs/internal/datanode/server"
	logger "github.com/cyb0225/gdfs/internal/pkg/log"
)

// load config file and start rpc server
func Run() {
	if err := config.NewConfig(); err != nil {
		log.Fatalf("read config failed: %s\n", err.Error())
	}
		
	logger.NewLogger(config.Cfg.Log)

	cfg := &server.ServerConfig{
		IP: config.Cfg.Addr.IP,
		Port: config.Cfg.Addr.Port,
		StoragePath: config.Cfg.StoragePath,
		NamenodeAddr: config.Cfg.NamenodeAddr,
		HeartBeatTime: config.Cfg.HeartBeatTime,
	}
	if err := server.RunServer(cfg); err != nil {
		log.Fatalf("start rpc server failed: %s\n", err.Error())
	}
}
