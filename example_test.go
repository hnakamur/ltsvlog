package ltsvlog_test

import (
	"errors"
	"os"

	"github.com/hnakamur/ltsvlog"
)

func ExampleNewLTSVLogger() {
	// Change the global logger to a logger which does not print level values.
	ltsvlog.Logger = ltsvlog.NewLTSVLogger(os.Stdout, true, ltsvlog.SetLevelLabel(""))
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func ExampleLTSVLogger_Debug() {
	if ltsvlog.Logger.DebugEnabled() {
		ltsvlog.Logger.Debug(ltsvlog.LV{"msg", "This is a debug message"},
			ltsvlog.LV{"key", "key1"}, ltsvlog.LV{"intValue", 234})
	}

	// Output example:
	// time:2017-05-19T20:39:45.112667Z	level:Debug	msg:This is a debug message	key:key1	intValue:234
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func ExampleLTSVLogger_Info() {
	ltsvlog.Logger.Info(ltsvlog.LV{"msg", "hello, world"}, ltsvlog.LV{"foo", "bar"},
		ltsvlog.LV{"nilValue", nil}, ltsvlog.LV{"bytes", []byte("a/b")})

	// Output example:
	// time:2017-05-19T20:42:23.631263Z	level:Info	msg:hello, world	foo:bar nilValue:<nil> bytes:0x612f62
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func ExampleLTSVLogger_Err() {
	b := func() error {
		return ltsvlog.Err(errors.New("some error")).Time().LV("key1", "value1").Stack()
	}
	a := func() error {
		return b()
	}
	err := a()
	if err != nil {
		ltsvlog.Logger.Err(err)
	}

	// Output example:
	// time:2017-05-19T20:45:56.597594Z	level:Error	err:some error	errtime:2017-05-20T05:45:56.597551Z	key1:value1	stack:[main.main.func1(0x1, 0xc42000e2a0) /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/err.go:11 +0x26b],[main.main.func2(0x4af4c0, 0xc42000e2a0) /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/err.go:14 +0x26],[main.main() /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/err.go:16 +0x65]
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}
