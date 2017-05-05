package ltsvlog_test

import (
	"errors"
	"os"

	"github.com/hnakamur/ltsvlog"
)

func Example() {
	if ltsvlog.Logger.DebugEnabled() {
		ltsvlog.Logger.Debug(ltsvlog.LV{"msg", "This is a debug message"},
			ltsvlog.LV{"key", "key1"}, ltsvlog.LV{"intValue", 234})
	}
	ltsvlog.Logger.Info(ltsvlog.LV{"msg", "hello, world"}, ltsvlog.LV{"key", "key1"},
		ltsvlog.LV{"value", "value1"})

	b := func() {
		err := errors.New("demo error")
		if err != nil {
			ltsvlog.Logger.ErrorWithStack(ltsvlog.LV{"err", err})
		}
	}
	a := func() {
		b()
	}
	a()

	ltsvlog.Logger.Info(ltsvlog.LV{"msg", "goodbye, world"}, ltsvlog.LV{"foo", "bar"},
		ltsvlog.LV{"nilValue", nil}, ltsvlog.LV{"bytes", []byte("a/b")})

	//Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func ExampleNewLTSVLogger() {
	// Change the global logger to a logger which does not print level values.
	ltsvlog.Logger = ltsvlog.NewLTSVLogger(os.Stdout, true, ltsvlog.SetLevelLabel(""))
}
