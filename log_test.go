package ltsvlog

import (
	"bytes"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/hnakamur/errstack"
)

func TestAppendTime(t *testing.T) {
	testCases := []struct {
		buf  []byte
		val  time.Time
		want string
	}{
		{buf: nil, val: time.Unix(0, 0).UTC(), want: "1970-01-01T00:00:00.000000Z"},
		{buf: nil, val: time.Date(2017, 5, 7, 22, 13, 59, 987654000, time.UTC), want: "2017-05-07T22:13:59.987654Z"},
		{buf: []byte("time:"), val: time.Date(2017, 5, 7, 22, 13, 59, 987654000, time.UTC), want: "time:2017-05-07T22:13:59.987654Z"},
	}
	for _, c := range testCases {
		buf := appendUTCTime(c.buf, c.val)
		want := []byte(c.want)
		if !bytes.Equal(buf, want) {
			t.Errorf("time value mismatch. got=%s, want=%s", string(buf), want)
		}
	}
}

func ExampleLTSVLogger_Info() {
	Logger = NewLTSVLogger(os.Stdout, false, UseLocalTimeZone(), SetLocalTimeZoneFormat("-0700"))
	Logger.Info().String("msg", "hello").Log()

	// Output example:
	// time:2019-10-21T22:05:06.784123+0900	level:Info	msg:hello

	// Actually we don't test the results.
	// This example is added just for document purpose.
}

func TestEventErr(t *testing.T) {
	var b bytes.Buffer
	l := NewLTSVLogger(&b, false, SetTimeLabel(""))
	l.Warn().Err(errstack.WithLV(errors.New("some error")).String("field1", "value1")).String("field2", "value2").Log()
	if got, want := b.String(), "level:Warn\terr:some error\tfield1:value1\tfield2:value2\n"; got != want {
		t.Errorf("result mismatch,\n got=%q,\nwant=%q", got, want)
	}
}
