package zap

import (
	"go.uber.org/zap"
)


type zaplogger struct {
	l *zap.Logger        // 不使用反射的方式打印, 性能稍高
	s *zap.SugaredLogger // 使用反射方式打印, 方便快捷
}

func NewZapLogger() *zaplogger {
	core := NewCore()

	logger := zap.New(core, zap.AddCaller())
	sugar := logger.Sugar()

	zl :=  &zaplogger{
		l: logger,
		s: sugar,
	}



	return zl
}

func (z *zaplogger) Debug(msg string, keyAndValues ...string) {

}

func (z *zaplogger) Info(msg string, keyAndValues ...string) {

}

func (z *zaplogger) Warn(msg string, keyAndValues ...string)
func (z *zaplogger) Error(msg string, keyAndValues ...string)
func (z *zaplogger) Panic(msg string, keyAndValues ...string)
func (z *zaplogger) Fatal(msg string, keyAndValues ...string)  

func (z *zaplogger) Debugf(format string, v ...interface{})
func (z *zaplogger) Infof(format string, v ...interface{})
func (z *zaplogger) Warnf(format string, v ...interface{})
func (z *zaplogger) Errorf(format string, v ...interface{})
func (z *zaplogger) Panicf(format string, v ...interface{})
func (z *zaplogger) Fatalf(format string, v ...interface{})
