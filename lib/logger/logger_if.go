package logger

type LoggerIf interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
	Err(err error)
	Panic(message string)
}
