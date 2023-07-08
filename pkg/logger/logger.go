package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
}

type logType struct {
	i, e  *log.Logger
	space string
}

func New(space string) Logger {
	return &logType{
		space: space,
		i:     log.New(os.Stdout, "I ", log.LstdFlags|log.Lmicroseconds),
		e:     log.New(os.Stdout, "E ", log.LstdFlags|log.Lmicroseconds),
	}
}

func (l *logType) Info(format string, args ...interface{}) {
	l.i.Printf("%v: %v", l.space, fmt.Sprintf(format, args...))
}
func (l *logType) Error(format string, args ...interface{}) {
	l.e.Printf("%v: %v", l.space, fmt.Sprintf(format, args...))
}
