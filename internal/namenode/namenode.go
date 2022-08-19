package namenode

import (
	"log"

	"github.com/cyb0225/gdfs/internal/namenode/config"
	"github.com/cyb0225/gdfs/internal/namenode/server"
	logger "github.com/cyb0225/gdfs/internal/pkg/log"
)

// load config file and start rpc server
func Run() {

	if err := config.NewConfig(); err != nil {
		log.Fatalf("read config failed: %s\n", err.Error())
	}

	// log.Printf("namenode config: %#v", config.Cfg)

	logger.NewLogger(config.Cfg.Log)

	cfg := &server.ServerConfig{
		Port: config.Cfg.Addr.Port,
		Backups: config.Cfg.BackupN,
		ChunkSize: config.Cfg.ChunkSize,
		StoragePath: config.Cfg.StoragePath,
		Expired: config.Cfg.Timeout,
	}
	
	if err := server.RunServer(cfg); err != nil {
		log.Fatalf("start rpc server failed: %s\n", err.Error())
	}
}