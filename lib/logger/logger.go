package logger

var loggerStr LoggerIf

func SetDefaultLogger() {
	loggerStr = NewLoggerZerolog(nil)
}

func Debug(message string) {
	loggerStr.Debug(message)
}

func Info(message string) {
	loggerStr.Info(message)
}

func Warn(message string) {
	loggerStr.Warn(message)
}

func Error(message string) {
	loggerStr.Error(message)
}

func Err(err error) {
	loggerStr.Err(err)
}

func Panic(message string) {
	loggerStr.Panic(message)
	panic("")
}
