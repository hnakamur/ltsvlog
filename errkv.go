package ltsvlog

import (
	"bytes"
	"runtime"
)

type ErrKV struct {
	error
	lvs []LV
}

func Err(err error) *ErrKV {
	return &ErrKV{
		error: err,
		lvs: []LV{
			{"err", err},
			{"stack", fullstack(2)},
		},
	}
}

func (e *ErrKV) KV(key string, value interface{}) *ErrKV {
	e.lvs = append(e.lvs, LV{key, value})
	return e
}

func (e *ErrKV) Error() string {
	return e.error.Error()
}

func (e *ErrKV) GetError() error {
	return e.error
}

func (e *ErrKV) ToLVs() []LV {
	return e.lvs
}

// fullstack formats a stack trace of the calling goroutine into buf
// in one line format which suitable for LTSV logs.
func fullstack(skip int) string {
	buf := make([]byte, 8192)
	var n int
	for {
		n = runtime.Stack(buf, false)
		if n < len(buf) {
			break
		}
		buf = make([]byte, len(buf)*2)
	}
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
	return string(p)
}
