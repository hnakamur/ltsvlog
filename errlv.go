package ltsvlog

import (
	"bytes"
	"runtime"
	"sync"
	"time"
)

var errLVPool = &sync.Pool{
	New: func() interface{} {
		return &ErrLV{
			buf: make([]byte, 8192),
		}
	},
}

// ErrLV is an error with label and value pairs.
//
// This is useful when you would like to log an error with
// additional labeled values later at the higher level of
// the callstack.
//
// ErrLV frees lower level functions from depending on loggers
// since ErrLV is just a data structure which holds
// an error, a stacktrace and labeld values.
//
// Please see the example at LTSVLogger.Err for an example usage.
type ErrLV struct {
	error
	buf []byte
}

// Err creates an ErrLV with the specified error.
func Err(err error) *ErrLV {
	e := errLVPool.Get().(*ErrLV)
	e.error = err
	e.buf = e.buf[:0]
	e.buf = append(e.buf, "err:"...)
	e.buf = appendValue(e.buf, err)
	return e
}

// Stack appends a stacktrace with label "stack" to ErrLV.
func (e *ErrLV) Stack() *ErrLV {
	e.buf = append(e.buf, "\tstack:"...)
	e.buf = appendStack(e.buf, 2)
	return e
}

// Time appends a current time with label "errtime" to ErrLV.
// This is useful when you log some time later after an error occurs.
func (e *ErrLV) Time() *ErrLV {
	e.buf = append(e.buf, "\terrtime:"...)
	e.buf = appendTime(e.buf, time.Now())
	return e
}

// LV appends a label value pair to ErrLV.
func (e *ErrLV) LV(key string, value interface{}) *ErrLV {
	e.buf = append(e.buf, '\t')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, ':')
	e.buf = appendValue(e.buf, value)
	return e
}

// LV returns the original error string without label and values appended.
func (e *ErrLV) Error() string {
	return string(e.buf)
}

// GetError returns the original error.
func (e *ErrLV) GetError() error {
	return e.error
}

var stackBufPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 8192)
	},
}

// appendStack appends a formated stack trace of the calling goroutine to buf
// in one line format which suitable for LTSV logs.
func appendStack(buf []byte, skip int) []byte {
	src := stackBufPool.Get().([]byte)
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
	stackBufPool.Put(src)
	return buf
}
