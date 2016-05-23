package ltsvlog

import (
	"bytes"
	"log"
	"testing"
)

func BenchmarkLTSVLog(b *testing.B) {
	w := new(bytes.Buffer)
	logger := NewLTSVLogger(w, false, nil)
	for i := 0; i < b.N; i++ {
		logger.Info(LV{"msg", "sample log message"}, LV{"key1", "value1"}, LV{"key2", "value2"})
	}
}

func BenchmarkStandardLog(b *testing.B) {
	w := new(bytes.Buffer)
	logger := log.New(w, "", log.LstdFlags|log.Lmicroseconds)
	for i := 0; i < b.N; i++ {
		logger.Printf("msg:sample log message\tkey1:%s\tkey2:%s", "value1", "value2")
	}
}
