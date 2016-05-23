package main

import (
	"os"

	"github.com/hnakamur/ltsvlog"
)

func main() {
	logger := ltsvlog.NewLTSVLogger(os.Stdout, true, nil)
	if logger.DebugEnabled() {
		logger.Debug(ltsvlog.LV{"msg", "This is a debug message"}, ltsvlog.LV{"key", "key1"}, ltsvlog.LV{"intValue", 234})
	}
	logger.Info(ltsvlog.LV{"msg", "hello, world"}, ltsvlog.LV{"key", "key1"}, ltsvlog.LV{"value", "value1"})
	logger.Info(ltsvlog.LV{"msg", "goodbye, world"}, ltsvlog.LV{"foo", "bar"}, ltsvlog.LV{"nilValue", nil}, ltsvlog.LV{"bytes", []byte("a/b")})
}
