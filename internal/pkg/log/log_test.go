package log

import (
	"testing"

	logger "github.com/cyb0225/gdfs/pkg/log"
)

func TestLog(t *testing.T) {
	logConfig := LogConfig{
		Module:     "info",
		LogPath: "./log/access.log",
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     3,
		Compress:   false,
	}

	NewLogger(&logConfig)
	for i := 0; i < 1; i++ {
		logger.Debug("debug", logger.String("string", "string"))
		logger.Info("info", logger.Int("int", 1))
		logger.Debugf("%s", "123")
	}
}
