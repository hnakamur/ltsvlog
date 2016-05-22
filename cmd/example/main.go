package main

import (
	"io"
	"log"
	"os"

	"github.com/hnakamur/ltsvlog"
)

type myLogger struct {
	logger *log.Logger
}

func newMyLogger(w io.Writer) *myLogger {
	return &myLogger{
		logger: log.New(w, "", 0),
	}
}

func (l *myLogger) Log(v interface{}) {
	l.logger.Print(v)
}

func main() {
	logger := ltsvlog.NewLTSVLogger(newMyLogger(os.Stdout), false)
	if logger.DebugEnabled() {
		logger.Debug(ltsvlog.LV{"msg", "This is a debug message"}, ltsvlog.LV{"key", "key1"}, ltsvlog.LV{"value", "value1"})
	}
	logger.Info(ltsvlog.LV{"msg", "hello, world"}, ltsvlog.LV{"key", "key1"}, ltsvlog.LV{"value", "value1"})
	logger.Info(ltsvlog.LV{"msg", "goodbye, world"}, ltsvlog.LV{"foo", "bar"})
}
