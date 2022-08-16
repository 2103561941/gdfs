package log

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// time format

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
