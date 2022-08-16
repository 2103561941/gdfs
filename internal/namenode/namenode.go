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

	if err := server.RunServer(); err != nil {
		log.Fatalf("start rpc server failed: %s\n", err.Error())
	}
}
