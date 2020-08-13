package logrus

import l "github.com/sirupsen/logrus"

type Logger struct {
    *l.Logger
}

func (t *Logger) ErrorStack(err interface{}, stack string) {
}

func NewLogger() *Logger {
    return &Logger{l.New()}
}
