package ltsvlog_test

import (
	"os"

	"github.com/hnakamur/errstack"
	ltsvlog "github.com/hnakamur/ltsvlog/v3"
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
		// In real usage, you might do a time consuming operation to
		// build values for logging, but it will be skipped when the debug
		// log is disabled.
		n := 234
		ltsvlog.Logger.Debug().String("msg", "This is a debug message").
			String("key", "key1").Int("intValue", n).Log()
	}

	// Output example:
	// time:2017-05-20T19:12:10.883958Z	level:Debug	msg:This is a debug message	key:key1	intValue:234
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func ExampleLTSVLogger_Info() {
	ltsvlog.Logger.Info().String("msg", "goodbye, world").String("foo", "bar").
		Fmt("nilValue", "%v", nil).HexBytes("bytes", []byte("a/b")).Log()

	// Output example:
	// time:2017-05-20T19:16:11.798840Z	level:Info	msg:goodbye, world	foo:bar	nilValue:<nil>	bytes:0x612f62
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func ExampleLTSVLoggerErr() {
	b := func() error {
		return errstack.New("some error")
	}
	a := func() error {
		if err := b(); err != nil {
			return errstack.Errorf("add some message here: %s", err)
		}
		return nil

	}
	err := a()
	if err != nil {
		ltsvlog.Logger.Err(err)
	}

	// Output example:
	// time:2019-10-21T22:00:34.549974Z	level:Error	err:add some message here: some error	reqID:req1	userID:1	stack:github.com/hnakamur/ltsvlog/v3_test.exampleErrInner@/home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example_err_test.go:23 github.com/hnakamur/ltsvlog/v3_test.exampleErrOuter@/home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example_err_test.go:16 github.com/hnakamur/ltsvlog/v3_test.ExampleLTSVLogger_Err@/home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example_err_test.go:9 testing.runExample@/usr/local/go/src/testing/run_example.go:62 testing.runExamples@/usr/local/go/src/testing/example.go:44 testing.(*M).Run@/usr/local/go/src/testing/testing.go:1118 main.main@_testmain.go:52 runtime.main@/usr/local/go/src/runtime/proc.go:203
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func ExampleEvent_String() {
	jsonStr := "{\n\t\"foo\": \"bar\\nbaz\"\n}\n"
	ltsvlog.Logger.Info().String("json", jsonStr).Log()

	// Output example:
	// time:2017-06-10T10:22:48.083226Z	level:Info	json:{\n\t"foo": "bar\\nbaz"\n}\n
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}
