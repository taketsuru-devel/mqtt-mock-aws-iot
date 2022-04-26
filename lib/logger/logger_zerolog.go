package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type loggerZerolog struct {
	logger zerolog.Logger
}

type LoggerZerologModFunc func(l *zerolog.Logger) *zerolog.Logger

func NewLoggerZerolog(modfunc LoggerZerologModFunc) *loggerZerolog {
	l := zerolog.New(os.Stdout)
	if modfunc != nil {
		ret := modfunc(&l)
		l = *ret
	}

	return &loggerZerolog{
		logger: l,
	}
}

func (l *loggerZerolog) Debug(message string) {
	l.logger.Debug().Msg(message)
}

func (l *loggerZerolog) Info(message string) {
	l.logger.Info().Msg(message)
}

func (l *loggerZerolog) Warn(message string) {
	l.logger.Warn().Msg(message)
}

func (l *loggerZerolog) Error(message string) {
	l.logger.Error().Msg(message)
}

func (l *loggerZerolog) Err(err error) {
	l.logger.Err(err).Msg("")
}

func (l *loggerZerolog) Panic(message string) {
	l.logger.Panic().Msg(message)
}
