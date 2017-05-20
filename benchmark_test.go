package ltsvlog_test

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

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
		logger.LV("msg", "hello").LV("key1", "value1").Info()
	}
}

func BenchmarkErrWithStackAndTime(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	run := func() error {
		return ltsvlog.Err(errors.New("some error")).Stack().Time().LV("key1", "value1")
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
		return ltsvlog.Err(errors.New("some error")).Stack().LV("key1", "value1")
	}

	logger := ltsvlog.NewLTSVLogger(tmpfile, false)
	for i := 0; i < b.N; i++ {
		err = run()
		logger.Err(err)
	}
}

func BenchmarkErrWithTime(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	run := func() error {
		return ltsvlog.Err(errors.New("some error")).Time().LV("key1", "value1")
	}

	logger := ltsvlog.NewLTSVLogger(tmpfile, false)
	for i := 0; i < b.N; i++ {
		err = run()
		logger.Err(err)
	}
}
