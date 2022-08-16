package log

func Debug(msg string, fields ...Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	logger.Warn(msg, fields...)
}

// it is different from field error
func Error(msg string, fields ...Field) {
	logger.Error(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	logger.Fatal(msg, fields...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v)
}
func Infof(format string, v ...interface{}) {
	logger.Infof(format, v)
}
func Warnf(format string, v ...interface{}) {
	logger.Warnf(format, v)
}
func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v)
}
func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v)
}
func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v)
}
