package client

import (
	"log"

	"github.com/cyb0225/gdfs/internal/client/cmd"
	"github.com/cyb0225/gdfs/internal/client/config"
	logger "github.com/cyb0225/gdfs/internal/pkg/log"
)

func Run() {
	if err := config.NewConfig(); err != nil {
		log.Fatalf("read config failed: %s\n", err.Error())
	}

	logger.NewLogger(config.Cfg.Log)

	cmd.Execute()
}
