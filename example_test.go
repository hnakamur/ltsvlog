package ltsvlog_test

import (
	"errors"
	"os"
	"time"

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
		n := 234
		ltsvlog.Logger.Debug().String("msg", "This is a debug message").
			String("key", "key1").Int64("intValue", int64(n)).Log()
	}

	// Output example:
	// time:2017-05-20T19:12:10.883958Z	level:Debug	msg:This is a debug message	key:key1	intValue:234
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func ExampleLTSVLogger_Info() {
	ltsvlog.Logger.Info().String("msg", "goodbye, world").String("foo", "bar").
		Sprintf("nilValue", "%v", nil).Hex("bytes", []byte("a/b")).Log()

	// Output example:
	// time:2017-05-20T19:16:11.798840Z	level:Info	msg:goodbye, world	foo:bar	nilValue:<nil>	bytes:0x612f62
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func ExampleLTSVLogger_Err() {
	b := func() error {
		return ltsvlog.Err(errors.New("some error")).UTCTime("errtime", time.Now()).String("key1", "value1").Stack("stack")
	}
	a := func() error {
		return b()
	}
	err := a()
	if err != nil {
		ltsvlog.Logger.Err(err)
	}

	// Output example:
	//time:2017-05-20T19:18:18.490029Z	level:Error	err:some error	errtime:2017-05-20T19:18:18.490007Z	key1:value1	stack:[main.b(0x0, 0x0) /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/err/main.go:22 +0x118],[main.a(0xc420016200, 0xc4200001a0) /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/err/main.go:18 +0x22],[main.main() /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/err/main.go:11 +0x22]
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}
