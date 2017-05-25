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
// If label is empty, "stack" is used.
func (e *ErrorEvent) Stack(label string) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	if label == "" {
		label = "stack"
	}
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

// Stringer appends a labeled string value to ErrorEvent.
// The value will be converted to a string with String() method.
func (e *ErrorEvent) Stringer(label string, value fmt.Stringer) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, escape(value.String())...)
	return e
}

// Byte appends a labeled byte value to ErrorEvent.
func (e *ErrorEvent) Byte(label string, value byte) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendHexByte(e.buf, value)
	return e
}

// Bytes appends a labeled bytes value to ErrorEvent.
func (e *ErrorEvent) Bytes(label string, value []byte) *ErrorEvent {
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

// Int appends a labeled int value to ErrorEvent.
func (e *ErrorEvent) Int(label string, value int) *ErrorEvent {
	return e.Int64(label, int64(value))
}

// Int8 appends a labeled int8 value to ErrorEvent.
func (e *ErrorEvent) Int8(label string, value int8) *ErrorEvent {
	return e.Int64(label, int64(value))
}

// Int16 appends a labeled int16 value to ErrorEvent.
func (e *ErrorEvent) Int16(label string, value int16) *ErrorEvent {
	return e.Int64(label, int64(value))
}

// Int32 appends a labeled int32 value to ErrorEvent.
func (e *ErrorEvent) Int32(label string, value int32) *ErrorEvent {
	return e.Int64(label, int64(value))
}

// Int64 appends a labeled int64 value to ErrorEvent.
func (e *ErrorEvent) Int64(label string, value int64) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendInt(e.buf, value, 10)
	return e
}

// Uint appends a labeled uint value to ErrorEvent.
func (e *ErrorEvent) Uint(label string, value uint) *ErrorEvent {
	return e.Uint64(label, uint64(value))
}

// Uint8 appends a labeled uint8 value to ErrorEvent.
func (e *ErrorEvent) Uint8(label string, value uint8) *ErrorEvent {
	return e.Uint64(label, uint64(value))
}

// Uint16 appends a labeled uint16 value to ErrorEvent.
func (e *ErrorEvent) Uint16(label string, value uint16) *ErrorEvent {
	return e.Uint64(label, uint64(value))
}

// Uint32 appends a labeled uint32 value to ErrorEvent.
func (e *ErrorEvent) Uint32(label string, value uint32) *ErrorEvent {
	return e.Uint64(label, uint64(value))
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

// Time appends a labeled formatted time value to ErrorEvent.
// The format is the same as that in the Go standard time package.
// If the format is empty, time.RFC3339 is used.
func (e *ErrorEvent) Time(label string, value time.Time, format string) *ErrorEvent {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	if format == "" {
		format = time.RFC3339
	}
	e.buf = append(e.buf, escape(value.Format(format))...)
	return e
}

// UTCTime appends a labeled time value to ErrorEvent.
// The time value is converted to UTC and then printed
// in the same format as the log time field.
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
