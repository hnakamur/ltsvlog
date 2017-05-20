package ltsvlog_test

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hnakamur/ltsvlog"
)

func BenchmarkInfo(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	logger := ltsvlog.NewLTSVLogger(tmpfile, false)
	for i := 0; i < b.N; i++ {
		logger.Info().String("msg", "hello").String("key1", "value1").Log()
	}
}

func BenchmarkInfoString(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	logger := ltsvlog.NewLTSVLogger(tmpfile, false)
	for i := 0; i < b.N; i++ {
		logger.Info().String("msg", "hello").String("key1", "value1").Log()
	}
}

func BenchmarkErrWithStackAndUTCTime(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	run := func() error {
		return ltsvlog.Err(errors.New("some error")).Stack("stack").UTCTime("errtime", time.Now()).String("key1", "value1")
	}

	logger := ltsvlog.NewLTSVLogger(tmpfile, false)
	for i := 0; i < b.N; i++ {
		err = run()
		logger.Err(err)
	}
}

func BenchmarkErrWithStack(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	run := func() error {
		return ltsvlog.Err(errors.New("some error")).Stack("stack").String("key1", "value1")
	}

	logger := ltsvlog.NewLTSVLogger(tmpfile, false)
	for i := 0; i < b.N; i++ {
		err = run()
		logger.Err(err)
	}
}

func BenchmarkErrWithUTCTime(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	run := func() error {
		return ltsvlog.Err(errors.New("some error")).UTCTime("errtime", time.Now()).String("key1", "value1")
	}

	logger := ltsvlog.NewLTSVLogger(tmpfile, false)
	for i := 0; i < b.N; i++ {
		err = run()
		logger.Err(err)
	}
}
