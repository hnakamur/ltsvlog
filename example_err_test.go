package ltsvlog_test

import (
	"github.com/hnakamur/errstack"
	ltsvlog "github.com/hnakamur/ltsvlog/v3"
)

func ExampleLTSVLogger_Err() {
	if err := exampleErrOuter(); err != nil {
		ltsvlog.Logger.Err(err)
	}
}

func exampleErrOuter() error {
	if err := exampleErrInner(); err != nil {
		return errstack.WithLV(errstack.Errorf("add some message here: %s", err), "userID", "user1")
	}
	return nil
}

func exampleErrInner() error {
	return errstack.WithLV(errstack.New("some error"), "reqID", "req1")
}
