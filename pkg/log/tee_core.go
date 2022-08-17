package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
Filename: filepath + filename (the place file stored)
MaxSize：the max mememory of a log file(MB)
MaxBackups：the max number of log files stored
MaxAges：the latest time stored old log files
Compress：whether to compress and archive old files
*/
type RotateOptions struct {
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type LevelEnablerFunc = zap.LevelEnablerFunc

type TeeOption struct {
	Filename string
	Ropt     RotateOptions
	Lef      LevelEnablerFunc
}

func newTeeWithRotateCore(tops []TeeOption) zapcore.Core {
	cores := []zapcore.Core{}
	cfg := zap.NewProductionConfig()

	cfg.EncoderConfig.EncodeTime = timeEncoder // add time format

	for _, top := range tops {

		lv := top.Lef

		w := &lumberjack.Logger{
			Filename:   top.Filename,
			MaxSize:    top.Ropt.MaxSize,
			MaxBackups: top.Ropt.MaxBackups,
			MaxAge:     top.Ropt.MaxAge,
			Compress:   top.Ropt.Compress,
		}

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.AddSync(w),
			lv,
		)

		cores = append(cores, core)
	}

	core := zapcore.NewTee(cores...)

	return core
}

func newTeeWithRotate(tops []TeeOption, level Level, opts ...Option) *Logger {

	core := newTeeWithRotateCore(tops)

	// if it is debug level, then I will start stderr output.
	if level == DebugLevel {
		stdCore := newCore(os.Stderr, DebugLevel)
		core = zapcore.NewTee(core, stdCore)
	}

	log := zap.New(core, opts...)
	suger := log.Sugar()

	logger := &Logger{
		l: log,
		s: suger,
	}

	return logger
}
