package log

import "github.com/cyb0225/gdfs/pkg/log/zap"

var (
	Log Logger
)

// Log 日志包需要对外实现的接口
type Logger interface {
	Debug(msg string, keyAndValues ...string)
	Info(msg string, keyAndValues ...string)
	Warn(msg string, keyAndValues ...string)
	Error(msg string, keyAndValues ...string)
	Panic(msg string, keyAndValues ...string)
	Fatal(msg string, keyAndValues ...string)

	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Panicf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

func Setup() error {
	Log = zap.NewZapLogger()
	return nil
}
