package log

import (
	"runtime"
	"strconv"

	"go.uber.org/zap"
)

var (
	logger *Logger
)

type Logger struct {
	l *zap.Logger        // 不使用反射的方式打印, 性能稍高
	s *zap.SugaredLogger // 使用反射方式打印, 方便快捷
}

// Example of new Logger
func ExNewLogger() {
	tops := []TeeOption{
		{
			Filename: "./storage/log/access.log",
			Ropt: RotateOptions{
				MaxSize:    1,
				MaxAge:     3,
				MaxBackups: 3,
				Compress:   false,
			},
			Lef: func(l Level) bool {
				return l <= InfoLevel
			},
		},
		{
			Filename: "./storage/log/error.log",
			Ropt: RotateOptions{
				MaxSize:    1,
				MaxAge:     3,
				MaxBackups: 3,
				Compress:   false,
			},
			Lef: func(l Level) bool {
				return l > InfoLevel
			},
		},
	}

	logger = NewTeeWithRotate(tops, DebugLevel)
}

func NewLogger(tops []TeeOption, level Level, opts ...Option) {
	logger = NewTeeWithRotate(tops, level, opts...)
}

func (z *Logger) Debug(msg string, fields ...Field) {
	fields = append(fields, getCurrentCaller()...)
	z.l.Debug(msg, fields...)
}

func (z *Logger) Info(msg string, fields ...Field) {
	fields = append(fields, getCurrentCaller()...)
	z.l.Info(msg, fields...)
}

func (z *Logger) Warn(msg string, fields ...Field) {
	fields = append(fields, getCurrentCaller()...)
	fields = append(fields, getFatherCaller()...)
	z.l.Warn(msg, fields...)
}

func (z *Logger) Error(msg string, fields ...Field) {
	fields = append(fields, getCurrentCaller()...)
	fields = append(fields, getFatherCaller()...)
	z.l.Error(msg, fields...)
}

func (z *Logger) Panic(msg string, fields ...Field) {
	fields = append(fields, getCurrentCaller()...)
	fields = append(fields, getFatherCaller()...)
	z.l.Panic(msg, fields...)
}

func (z *Logger) Fatal(msg string, fields ...Field) {
	fields = append(fields, getCurrentCaller()...)
	fields = append(fields, getFatherCaller()...)
	z.l.Fatal(msg, fields...)
}

func (z *Logger) Debugf(format string, v ...interface{}) {
	z.s.Debugf(format, v)
}
func (z *Logger) Infof(format string, v ...interface{}) {
	z.s.Infof(format, v)
}
func (z *Logger) Warnf(format string, v ...interface{}) {
	z.s.Warnf(format, v)
}
func (z *Logger) Errorf(format string, v ...interface{}) {
	z.s.Errorf(format, v)
}
func (z *Logger) Panicf(format string, v ...interface{}) {
	z.s.Panicf(format, v)
}
func (z *Logger) Fatalf(format string, v ...interface{}) {
	z.s.Fatalf(format, v)
}

const (
	currentFunc = 3
	fatherFunc  = 4
)

// 3 equals the current function, if the number is bigger, then you can get who use this function
func getCallerInfoForLog(num int, field string) (callerFields []zap.Field) {
	_, file, line, ok := runtime.Caller(num)
	if !ok {
		return
	}

	file = file + ":" + strconv.Itoa(line)

	callerFields = append(callerFields, zap.String(field, file))
	return
}

func getCurrentCaller() (callerFields []zap.Field) {
	return getCallerInfoForLog(currentFunc+1, "currentFunc")
}

func getFatherCaller() (callerFields []zap.Field) {
	return getCallerInfoForLog(fatherFunc+1, "fatherFunc")
}
