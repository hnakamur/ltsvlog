package ltsvlog_test

import (
	"errors"
	"fmt"
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
		Sprintf("nilValue", "%v", nil).Bytes("bytes", []byte("a/b")).Log()

	// Output example:
	// time:2017-05-20T19:16:11.798840Z	level:Info	msg:goodbye, world	foo:bar	nilValue:<nil>	bytes:0x612f62
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func ExampleLTSVLogger_Err() {
	b := func() error {
		return ltsvlog.Err(errors.New("some error")).String("key1", "value1").Stack("")
	}
	a := func() error {
		err := b()
		if err != nil {
			return ltsvlog.WrapErr(err, func(err error) error {
				return fmt.Errorf("add explanation here, err=%v", err)
			})
		}
		return nil

	}
	err := a()
	if err != nil {
		ltsvlog.Logger.Err(err)
	}

	// Output example:
	// time:2017-06-01T16:52:33.959914Z	level:Error	err:add explanation here, err=some error	key1:value1	stack:[main.b(0x2000, 0x38) /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/main.go:35 +0xc8],[main.a(0xc42001a240, 0x4c6e2e) /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/main.go:25 +0x22],[main.main() /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/main.go:18 +0x11e]
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func ExampleEvent_String() {
	jsonStr := "{\n\t\"foo\": \"bar\\nbaz\"\n}\n"
	ltsvlog.Logger.Info().String("json", jsonStr).Log()

	// Output example:
	// time:2017-06-10T10:22:48.083226Z        level:Info      json:{\n\t"foo": "bar\\nbaz"\n}\n
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}
