package ltsvlog

import (
	"bytes"
	"runtime"
	"time"
)

type ErrLV interface {
	LV(key string, value interface{}) ErrLV
	Stack() ErrLV
	Time() ErrLV

	Error() string
	GetError() error
	ToLVs() []LV
}

type errLV struct {
	error
	lvs []LV
}

func Err(err error) ErrLV {
	e := &errLV{
		error: err,
		lvs:   make([]LV, 0, 8),
	}
	e.lvs = append(e.lvs, LV{"err", err})
	return e
}

func (e *errLV) Stack() ErrLV {
	e.lvs = append(e.lvs, LV{"stack", fullstack(2)})
	return e
}

func (e *errLV) Time() ErrLV {
	e.lvs = append(e.lvs, LV{"errtime", formatTime(time.Now())})
	return e
}

func (e *errLV) LV(key string, value interface{}) ErrLV {
	e.lvs = append(e.lvs, LV{key, value})
	return e
}

func (e *errLV) Error() string {
	return e.error.Error()
}

func (e *errLV) GetError() error {
	return e.error
}

func (e *errLV) ToLVs() []LV {
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

func formatTime(t time.Time) string {
	buf := []byte("0000-00-00T00:00:00.000000Z")
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	itoa(buf[:4], year, 4)
	itoa(buf[5:7], int(month), 2)
	itoa(buf[8:10], day, 2)
	itoa(buf[11:13], hour, 2)
	itoa(buf[14:16], min, 2)
	itoa(buf[17:19], sec, 2)
	itoa(buf[20:26], t.Nanosecond()/1e3, 6)
	return string(buf)
}
