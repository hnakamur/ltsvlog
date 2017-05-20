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
			lvs: make([]LV, 0, 8),
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
	lvs []LV
}

// Err creates an ErrLV with the specified error.
func Err(err error) *ErrLV {
	e := errLVPool.Get().(*ErrLV)
	e.error = err
	e.lvs = e.lvs[:0]
	e.lvs = append(e.lvs, LV{"err", err})
	return e
}

// Stack appends a stacktrace with label "stack" to ErrLV.
func (e *ErrLV) Stack() *ErrLV {
	e.lvs = append(e.lvs, LV{"stack", fullstack(2)})
	return e
}

// Time appends a current time with label "errtime" to ErrLV.
// This is useful when you log some time later after an error occurs.
func (e *ErrLV) Time() *ErrLV {
	e.lvs = append(e.lvs, LV{"errtime", formatTime(time.Now())})
	return e
}

// LV appends a label value pair to ErrLV.
func (e *ErrLV) LV(key string, value interface{}) *ErrLV {
	e.lvs = append(e.lvs, LV{key, value})
	return e
}

// LV returns the original error string without label and values appended.
func (e *ErrLV) Error() string {
	return e.error.Error()
}

// GetError returns the original error.
func (e *ErrLV) GetError() error {
	return e.error
}

// toLVs converts a ErrLV to a LV slice which can be passed to ltsv.LogWriter.Error.
func (e *ErrLV) toLVs() []LV {
	return e.lvs
}

var stackBufPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 8192)
	},
}

// fullstack formats a stack trace of the calling goroutine into buf
// in one line format which suitable for LTSV logs.
func fullstack(skip int) string {
	buf := stackBufPool.Get().([]byte)
	var n int
	for {
		n = runtime.Stack(buf, false)
		if n < len(buf) {
			break
		}
		buf = make([]byte, len(buf)*2)
	}
	bufToPut := buf
	buf = buf[:n]

	// NOTE: We reuse the same buffer here.
	p := buf[:0]

	for j := 0; j < 1+2*skip; j++ {
		i := bytes.IndexByte(buf, '\n')
		buf = buf[i+1:]
	}

	for len(buf) > 0 {
		p = append(p, '[')
		i := bytes.IndexByte(buf, '\n')
		p = append(p, buf[:i]...)
		p = append(p, ' ')
		buf = buf[i+2:]
		i = bytes.IndexByte(buf, '\n')
		p = append(p, buf[:i]...)
		p = append(p, ']')
		buf = buf[i+1:]
		if len(buf) > 0 {
			p = append(p, ',')
		}
	}
	s := string(p)
	stackBufPool.Put(bufToPut)
	return s
}
