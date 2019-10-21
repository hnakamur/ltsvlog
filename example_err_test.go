package ltsvlog_test

import (
	"github.com/hnakamur/errstack"
	ltsvlog "github.com/hnakamur/ltsvlog/v3"
)

func ExampleLTSVLogger_Err() {
	if err := exampleErrOuter(); err != nil {
		ltsvlog.Logger.Err(err)
	}

	// Output example:
	// time:2019-10-21T22:05:06.784512Z	level:Error	err:add some message here: some error	reqID:req1	userID:1	stack:github.com/hnakamur/ltsvlog/v3_test.exampleErrInner@/home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example_err_test.go:24 github.com/hnakamur/ltsvlog/v3_test.exampleErrOuter@/home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example_err_test.go:17 github.com/hnakamur/ltsvlog/v3_test.ExampleLTSVLogger_Err@/home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example_err_test.go:9 testing.runExample@/usr/local/go/src/testing/run_example.go:62 testing.runExamples@/usr/local/go/src/testing/example.go:44 testing.(*M).Run@/usr/local/go/src/testing/testing.go:1118 main.main@_testmain.go:52 runtime.main@/usr/local/go/src/runtime/proc.go:203
	// Output:

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func exampleErrOuter() error {
	if err := exampleErrInner(); err != nil {
		return errstack.WithLV(errstack.Errorf("add some message here: %s", err)).Int64("userID", 1)
	}
	return nil
}

func exampleErrInner() error {
	return errstack.WithLV(errstack.New("some error")).String("reqID", "req1")
}
