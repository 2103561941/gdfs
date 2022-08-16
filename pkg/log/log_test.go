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
	Debug("debug")
	Error("error")
}
