// core 设计

package zap

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

func NewCore() zapcore.Core {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	level := getLogLevel("debug")

	core := zapcore.NewCore(encoder, writeSyncer, level)

	return core
}


func getLogLevel(opt string) (level zapcore.Level) {

	switch opt {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}
	return level
}

/*
Filename: filepath + filename (the place file stored)
MaxSize：the max mememory of a log file(MB)
MaxBackups：the max number of log files stored
MaxAges：the latest time stored old log files
Compress：whether to compress and archive old files
*/

// file save and split
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./storage/log/debug.log", //
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
