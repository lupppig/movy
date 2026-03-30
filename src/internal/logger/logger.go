package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	log zerolog.Logger
}



func NewLogger() *Logger {
	 fileLogger := &lumberjack.Logger{
		Filename:   "logs/app.log", // make sure logs/ exists
		MaxSize:    10,             
		MaxBackups: 5,              
		MaxAge:     7,              
		Compress:   true,
	}

	multi := io.MultiWriter(os.Stdout, fileLogger)
	logger := zerolog.New(multi).With().Timestamp().Logger()
	return &Logger{log: logger}
}


func (l *Logger) Info() *zerolog.Event{
	event := l.log.Info()
	return event
}


func (l *Logger) Debug() *zerolog.Event {
	event := l.log.Debug()
	return event
}


func (l *Logger) Warn() *zerolog.Event {
	event := l.log.Warn()
	return event
}


func (l *Logger) Error() *zerolog.Event {
	event := l.log.Error()
	return event
}

func (l *Logger) Fatal() *zerolog.Event {
	event := l.log.Fatal()
	return event
}