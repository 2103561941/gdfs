// 设置打印格式

package zap

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 整体封装
func getEncoder() zapcore.Encoder {

	encoderConfig := zap.NewProductionEncoderConfig()
	zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 可以支持写入不同的文件
	zapcore.NewTee()

	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 时间格式
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}

// Console 格式

// Json 格式
