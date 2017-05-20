package ltsvlog

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var errLVsPool = &sync.Pool{
	New: func() interface{} {
		return &ErrLVs{
			buf: make([]byte, 8192),
		}
	},
}

// ErrLVs is an error with label and value pairs.
//
// This is useful when you would like to log an error with
// additional labeled values later at the higher level of
// the callstack.
//
// ErrLVs frees lower level functions from depending on loggers
// since ErrLVs is just a data structure which holds
// an error, a stacktrace and labeld values.
//
// Please see the example at LTSVLogger.Err for an example usage.
type ErrLVs struct {
	error
	buf []byte
}

// Err creates an ErrLVs with the specified error.
func Err(err error) *ErrLVs {
	e := errLVsPool.Get().(*ErrLVs)
	e.error = err
	e.buf = e.buf[:0]
	e.buf = append(e.buf, "err:"...)
	e.buf = append(e.buf, escape(fmt.Sprintf("%+v", err))...)
	return e
}

// Stack appends a stacktrace with label "stack" to ErrLVs.
func (e *ErrLVs) Stack(label string) *ErrLVs {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendStack(e.buf, 2)
	return e
}

func (e *ErrLVs) String(label string, value string) *ErrLVs {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, escape(value)...)
	return e
}

func (e *ErrLVs) Hex(label string, value []byte) *ErrLVs {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendHexBytes(e.buf, value)
	return e
}

func (e *ErrLVs) Sprintf(label, format string, a ...interface{}) *ErrLVs {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, escape(fmt.Sprintf(format, a...))...)
	return e
}

func (e *ErrLVs) Bool(label string, value bool) *ErrLVs {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendBool(e.buf, value)
	return e
}

func (e *ErrLVs) Int64(label string, value int64) *ErrLVs {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendInt(e.buf, value, 10)
	return e
}

func (e *ErrLVs) Uint64(label string, value uint64) *ErrLVs {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendUint(e.buf, value, 10)
	return e
}

func (e *ErrLVs) Float32(label string, value float32) *ErrLVs {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, strconv.FormatFloat(float64(value), 'g', -1, 32)...)
	return e
}

func (e *ErrLVs) Float64(label string, value float64) *ErrLVs {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, strconv.FormatFloat(value, 'g', -1, 64)...)
	return e
}

func (e *ErrLVs) UTCTime(label string, value time.Time) *ErrLVs {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendUTCTime(e.buf, value)
	return e
}

// LV returns the original error string without label and values appended.
func (e *ErrLVs) Error() string {
	return string(e.buf)
}

// GetError returns the original error.
func (e *ErrLVs) GetError() error {
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
