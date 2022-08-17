package log

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestZap(t *testing.T) {
	ExNewLogger()
	err := os.RemoveAll("./storage/log")
	if err != nil {
		Debug("", zap.Error(err))
	}

	for i := 0; i < 1000000; i++ {
		Debug("debug", Int("debug", i))
	}
}
