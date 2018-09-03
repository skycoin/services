package logger

import "github.com/sirupsen/logrus"

type Logger struct {
	*logrus.Entry
}

func NewLogger(name string) *Logger {
	logger := &Logger{
		Entry: logrus.WithField("service",name),
	}

	return logger
}
