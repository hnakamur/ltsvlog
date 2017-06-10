package ltsvlog

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Error is an error with label and value pairs.
// *Error implements the error interface so you can
// return *Error as an error.
//
// This is useful when you would like to log an error with
// additional labeled values later at the higher level of
// the callstack.
//
// Error frees lower level functions from depending on loggers
// since Error is just a data structure which holds
// an error, a stacktrace and labeled values.
//
// Please see the example at LTSVLogger.Err for an example usage.
type Error struct {
	error
	originalErr error
	buf         []byte
}

// Err creates an Error with the specified error.
func Err(err error) *Error {
	return &Error{
		error:       err,
		originalErr: err,
		buf:         make([]byte, 0, 8192),
	}
}

// WrapErr wraps an Error or a plain error and returns a new error.
func WrapErr(err error, wrapper func(err error) error) *Error {
	e, ok := err.(*Error)
	if !ok {
		e = Err(err)
	}

	if wrapper != nil {
		e.error = wrapper(e.error)
	}
	return e
}

// Stack appends a stacktrace with label "stack" to Error.
// If label is empty, "stack" is used.
func (e *Error) Stack(label string) *Error {
	e.buf = append(e.buf, '\t')
	if label == "" {
		label = "stack"
	}
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendStack(e.buf, 2)
	return e
}

// String appends a labeled string value to Error.
func (e *Error) String(label string, value string) *Error {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, escape(value)...)
	return e
}

// Stringer appends a labeled string value to Error.
// The value will be converted to a string with String() method.
func (e *Error) Stringer(label string, value fmt.Stringer) *Error {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, escape(value.String())...)
	return e
}

// Byte appends a labeled byte value to Error.
func (e *Error) Byte(label string, value byte) *Error {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendHexByte(e.buf, value)
	return e
}

// Bytes appends a labeled bytes value in hex format to Error.
func (e *Error) Bytes(label string, value []byte) *Error {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendHexBytes(e.buf, value)
	return e
}

// Sprintf appends a labeled formatted string value to Error.
func (e *Error) Sprintf(label, format string, a ...interface{}) *Error {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, escape(fmt.Sprintf(format, a...))...)
	return e
}

// Bool appends a labeled bool value to Error.
func (e *Error) Bool(label string, value bool) *Error {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendBool(e.buf, value)
	return e
}

// Int appends a labeled int value to Error.
func (e *Error) Int(label string, value int) *Error {
	return e.Int64(label, int64(value))
}

// Int8 appends a labeled int8 value to Error.
func (e *Error) Int8(label string, value int8) *Error {
	return e.Int64(label, int64(value))
}

// Int16 appends a labeled int16 value to Error.
func (e *Error) Int16(label string, value int16) *Error {
	return e.Int64(label, int64(value))
}

// Int32 appends a labeled int32 value to Error.
func (e *Error) Int32(label string, value int32) *Error {
	return e.Int64(label, int64(value))
}

// Int64 appends a labeled int64 value to Error.
func (e *Error) Int64(label string, value int64) *Error {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendInt(e.buf, value, 10)
	return e
}

// Uint appends a labeled uint value to Error.
func (e *Error) Uint(label string, value uint) *Error {
	return e.Uint64(label, uint64(value))
}

// Uint8 appends a labeled uint8 value to Error.
func (e *Error) Uint8(label string, value uint8) *Error {
	return e.Uint64(label, uint64(value))
}

// Uint16 appends a labeled uint16 value to Error.
func (e *Error) Uint16(label string, value uint16) *Error {
	return e.Uint64(label, uint64(value))
}

// Uint32 appends a labeled uint32 value to Error.
func (e *Error) Uint32(label string, value uint32) *Error {
	return e.Uint64(label, uint64(value))
}

// Uint64 appends a labeled uint64 value to Error.
func (e *Error) Uint64(label string, value uint64) *Error {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendUint(e.buf, value, 10)
	return e
}

// Float32 appends a labeled float32 value to Error.
func (e *Error) Float32(label string, value float32) *Error {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, strconv.FormatFloat(float64(value), 'g', -1, 32)...)
	return e
}

// Float64 appends a labeled float64 value to Error.
func (e *Error) Float64(label string, value float64) *Error {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, strconv.FormatFloat(value, 'g', -1, 64)...)
	return e
}

// Time appends a labeled formatted time value to Error.
// The format is the same as that in the Go standard time package.
// If the format is empty, time.RFC3339 is used.
func (e *Error) Time(label string, value time.Time, format string) *Error {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	if format == "" {
		format = time.RFC3339
	}
	e.buf = append(e.buf, escape(value.Format(format))...)
	return e
}

// UTCTime appends a labeled time value to Error.
// The time value is converted to UTC and then printed
// in the same format as the log time field, that is
// the ISO8601 format with microsecond precision and
// the timezone "Z".
func (e *Error) UTCTime(label string, value time.Time) *Error {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendUTCTime(e.buf, value)
	return e
}

// Error returns the error string without labeled values.
func (e *Error) Error() string {
	return e.error.Error()
}

// Format formats the error. With %+v and %+q, labeled values are
// appended to the error message in LTSV format.
func (e *Error) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		if s.Flag('+') {
			buf := make([]byte, 0, 8192)
			buf = e.AppendErrorWithValues(buf)
			s.Write(buf)
		} else {
			io.WriteString(s, e.Error())
		}
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		if s.Flag('+') {
			buf := make([]byte, 0, 8192)
			buf = e.AppendErrorWithValues(buf)
			fmt.Fprintf(s, "%q", buf)
		} else {
			fmt.Fprintf(s, "%q", e.Error())
		}
	}
}

// AppendErrorWithValues appends the error string with labeled values to a byte buffer.
func (e *Error) AppendErrorWithValues(buf []byte) []byte {
	buf = append(buf, "err:"...)
	buf = append(buf, escape(e.Error())...)
	return append(buf, e.buf...)
}

// OriginalError returns the original error.
func (e *Error) OriginalError() error {
	return e.originalErr
}

// appendStack appends a formated stack trace of the calling goroutine to buf
// in one line format which suitable for LTSV logs.
func appendStack(buf []byte, skip int) []byte {
	src := bufPool.Get().([]byte)
	var n int
	for {
		n = runtime.Stack(src, false)
		if n < len(src) {
			break
		}
		src = make([]byte, len(src)*2)
	}

	p := src[:n]
	for j := 0; j < 1+2*skip; j++ {
		i := bytes.IndexByte(p, '\n')
		p = p[i+1:]
	}

	for len(p) > 0 {
		buf = append(buf, '[')
		i := bytes.IndexByte(p, '\n')
		buf = append(buf, p[:i]...)
		buf = append(buf, ' ')
		p = p[i+2:]
		i = bytes.IndexByte(p, '\n')
		buf = append(buf, p[:i]...)
		buf = append(buf, ']')
		p = p[i+1:]
		if len(p) > 0 {
			buf = append(buf, ',')
		}
	}
	bufPool.Put(src)
	return buf
}

var bufPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 8192)
	},
}
