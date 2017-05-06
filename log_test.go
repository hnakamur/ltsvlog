package ltsvlog

import (
	"bytes"
	"math"
	"testing"
	"time"
)

func TestAppendValueNil(t *testing.T) {
	buf := appendValue(nil, nil)
	want := []byte("<nil>")
	if !bytes.Equal(buf, want) {
		t.Errorf("nil value mismatch. got=%s, want=%s", string(buf), string(want))
	}
}

func TestAppendValueString(t *testing.T) {
	val := "string"
	buf := appendValue(nil, val)
	want := []byte("string")
	if !bytes.Equal(buf, want) {
		t.Errorf("string value mismatch. got=%s, want=%s", string(buf), string(want))
	}
}

func TestAppendValueUint8(t *testing.T) {
	var val uint8 = math.MaxUint8
	buf := appendValue(nil, val)
	want := []byte("255")
	if !bytes.Equal(buf, want) {
		t.Errorf("uint8 value mismatch. got=%s, want=%s", string(buf), string(want))
	}
}

func TestAppendValueUint16(t *testing.T) {
	var val uint16 = math.MaxUint16
	buf := appendValue(nil, val)
	want := []byte("65535")
	if !bytes.Equal(buf, want) {
		t.Errorf("uint16 value mismatch. got=%s, want=%s", string(buf), string(want))
	}
}

func TestAppendValueUint32(t *testing.T) {
	var val uint32 = math.MaxUint32
	buf := appendValue(nil, val)
	want := []byte("4294967295")
	if !bytes.Equal(buf, want) {
		t.Errorf("uint32 value mismatch. got=%s, want=%s", string(buf), string(want))
	}
}

func TestAppendValueUint64(t *testing.T) {
	var val uint64 = math.MaxUint64
	buf := appendValue(nil, val)
	want := []byte("18446744073709551615")
	if !bytes.Equal(buf, want) {
		t.Errorf("uint64 value mismatch. got=%s, want=%s", string(buf), string(want))
	}
}

func TestAppendValueInt8(t *testing.T) {
	var val int8 = math.MaxInt8
	buf := appendValue(nil, val)
	want := []byte("127")
	if !bytes.Equal(buf, want) {
		t.Errorf("int8 value mismatch. got=%s, want=%s", string(buf), string(want))
	}
}

func TestAppendValueInt16(t *testing.T) {
	var val int16 = math.MaxInt16
	buf := appendValue(nil, val)
	want := []byte("32767")
	if !bytes.Equal(buf, want) {
		t.Errorf("int16 value mismatch. got=%s, want=%s", string(buf), string(want))
	}
}

func TestAppendValueInt32(t *testing.T) {
	var val int32 = math.MaxInt32
	buf := appendValue(nil, val)
	want := []byte("2147483647")
	if !bytes.Equal(buf, want) {
		t.Errorf("int32 value mismatch. got=%s, want=%s", string(buf), string(want))
	}
}

func TestAppendValueInt64(t *testing.T) {
	var val int64 = math.MaxInt64
	buf := appendValue(nil, val)
	want := []byte("9223372036854775807")
	if !bytes.Equal(buf, want) {
		t.Errorf("int64 value mismatch. got=%s, want=%s", string(buf), string(want))
	}
}

func TestAppendValueBool(t *testing.T) {
	testCases := []struct {
		val  bool
		want string
	}{
		{val: true, want: "true"},
		{val: false, want: "false"},
	}
	for _, c := range testCases {
		buf := appendValue(nil, c.val)
		want := []byte(c.want)
		if !bytes.Equal(buf, want) {
			t.Errorf("bool value mismatch. got=%s, want=%s", string(buf), string(want))
		}
	}
}

func TestAppendValueFloat32(t *testing.T) {
	testCases := []struct {
		val  float32
		want string
	}{
		{val: 0, want: "0"},
		{val: math.SmallestNonzeroFloat32, want: "1e-45"},
		{val: math.MaxFloat32, want: "3.4028235e+38"},
		{val: -math.SmallestNonzeroFloat32, want: "-1e-45"},
		{val: -math.MaxFloat32, want: "-3.4028235e+38"},
		{val: float32(math.NaN()), want: "NaN"},
	}
	for _, c := range testCases {
		buf := appendValue(nil, c.val)
		want := []byte(c.want)
		if !bytes.Equal(buf, want) {
			t.Errorf("float32 value mismatch. got=%s, want=%s", string(buf), string(want))
		}
	}
}

func TestAppendValueFloat64(t *testing.T) {
	testCases := []struct {
		val  float64
		want string
	}{
		{val: 0, want: "0"},
		{val: math.SmallestNonzeroFloat64, want: "5e-324"},
		{val: math.MaxFloat64, want: "1.7976931348623157e+308"},
		{val: -math.SmallestNonzeroFloat64, want: "-5e-324"},
		{val: -math.MaxFloat64, want: "-1.7976931348623157e+308"},
		{val: float64(math.NaN()), want: "NaN"},
	}
	for _, c := range testCases {
		buf := appendValue(nil, c.val)
		want := []byte(c.want)
		if !bytes.Equal(buf, want) {
			t.Errorf("float64 value mismatch. got=%s, want=%s", string(buf), string(want))
		}
	}
}

func TestAppendTime(t *testing.T) {
	testCases := []struct {
		val  time.Time
		want string
	}{
		{val: time.Unix(0, 0).UTC(), want: "1970-01-01T00:00:00.000000Z"},
		{val: time.Date(2017, 5, 7, 22, 13, 59, 987654000, time.UTC), want: "2017-05-07T22:13:59.987654Z"},
	}
	for _, c := range testCases {
		buf := appendTime(nil, c.val)
		want := []byte(c.want)
		if !bytes.Equal(buf, want) {
			t.Errorf("time value mismatch. got=%s, want=%s", string(buf), want)
		}
	}
}
