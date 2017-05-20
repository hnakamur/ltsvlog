package ltsvlog

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var errorEventPool = &sync.Pool{
	New: func() interface{} {
		return &ErrorEvent{
			buf: make([]byte, 8192),
		}
	},
}

// ErrorEvent is an error with label and value pairs.
// *ErrroEvent implements the error interface so you can
// return *ErrorEvent as an error.
//
// This is useful when you would like to log an error with
// additional labeled values later at the higher level of
// the callstack.
//
// ErrorEvent frees lower level functions from depending on loggers
// since ErrorEvent is just a data structure which holds
// an error, a stacktrace and labeled values.
//
// However LTSVLogger.Err depends on the internal structure
// of ErrorEvent, so you cannnot use ErrorEvent with the
// other logging library than this package.
//
// Please see the example at LTSVLogger.Err for an example usage.
type ErrorEvent struct {
	error
	buf []byte
}

// Err creates an ErrorEvent with the specified error.
func Err(err error) *ErrorEvent {
	e := errorEventPool.Get().(*ErrorEvent)
	e.error = err
	e.buf = e.buf[:0]
	e.buf = append(e.buf, "err:"...)
	e.buf = append(e.buf, escape(err.Error())...)
	return e
}

// Stack appends a stacktrace with label "stack" to ErrorEvent.
func (e *ErrorEvent) Stack(label string) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendStack(e.buf, 2)
	return e
}

// String appends a labeled string value to ErrorEvent.
func (e *ErrorEvent) String(label string, value string) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, escape(value)...)
	return e
}

// Hex appends a labeled hex value to ErrorEvent.
func (e *ErrorEvent) Hex(label string, value []byte) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendHexBytes(e.buf, value)
	return e
}

// Sprintf appends a labeled formatted string value to ErrorEvent.
func (e *ErrorEvent) Sprintf(label, format string, a ...interface{}) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, escape(fmt.Sprintf(format, a...))...)
	return e
}

// Bool appends a labeled bool value to ErrorEvent.
func (e *ErrorEvent) Bool(label string, value bool) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendBool(e.buf, value)
	return e
}

// Int64 appends a labeled int64 value to ErrorEvent.
func (e *ErrorEvent) Int64(label string, value int64) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendInt(e.buf, value, 10)
	return e
}

// Uint64 appends a labeled uint64 value to ErrorEvent.
func (e *ErrorEvent) Uint64(label string, value uint64) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendUint(e.buf, value, 10)
	return e
}

// Float32 appends a labeled float32 value to ErrorEvent.
func (e *ErrorEvent) Float32(label string, value float32) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, strconv.FormatFloat(float64(value), 'g', -1, 32)...)
	return e
}

// Float64 appends a labeled float64 value to ErrorEvent.
func (e *ErrorEvent) Float64(label string, value float64) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, strconv.FormatFloat(value, 'g', -1, 64)...)
	return e
}

// UTCTime appends a labeled UTC time value to ErrorEvent.
func (e *ErrorEvent) UTCTime(label string, value time.Time) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendUTCTime(e.buf, value)
	return e
}

// Error returns the error string with labeled values in the LTSV format.
func (e *ErrorEvent) Error() string {
	return string(e.buf)
}

// OriginalError returns the original error.
func (e *ErrorEvent) OriginalError() error {
	return e.error
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
