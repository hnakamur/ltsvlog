package ltsvlog_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/hnakamur/ltsvlog"
	ltsv "github.com/hnakamur/zap-ltsv"
	"go.uber.org/zap"
)

func BenchmarkLTSVLog(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	logger := ltsvlog.NewLTSVLogger(tmpfile, false)
	for i := 0; i < b.N; i++ {
		logger.Info(ltsvlog.LV{"msg", "sample log message"}, ltsvlog.LV{"key1", "value1"}, ltsvlog.LV{"key2", "value2"})
	}
}

// NOTE: This does not produce a proper LTSV log since a log record does not have the "time: label.
// This is used just for benchmark comparison.
func BenchmarkStandardLog(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	logger := log.New(tmpfile, "", log.LstdFlags|log.Lmicroseconds)
	for i := 0; i < b.N; i++ {
		logger.Printf("level:Info\tmsg:sample log message\tkey1:%s\tkey2:%s", "value1", "value2")
	}
}

func init() {
	err := ltsv.RegisterLTSVEncoder()
	if err != nil {
		panic(err)
	}
}

func BenchmarkZapLTSVProductionLog(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	cfg := ltsv.NewProductionConfig()
	cfg.OutputPaths = []string{tmpfile.Name()}
	logger, err := cfg.Build()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		logger.Info("sample log message", zap.String("key1", "value1"), zap.String("key2", "value2"))
	}
}

func BenchmarkZapLTSVDevelopmentLog(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	cfg := ltsv.NewDevelopmentConfig()
	cfg.OutputPaths = []string{tmpfile.Name()}
	logger, err := cfg.Build()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		logger.Info("sample log message", zap.String("key1", "value1"), zap.String("key2", "value2"))
	}
}
