package ltsvlog

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{
			buf: make([]byte, 8192),
		}
	},
}

type Event struct {
	logger  *LTSVLogger
	enabled bool
	buf     []byte
}

func (e *Event) String(label string, value string) *Event {
	if !e.enabled {
		return e
	}
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, escape(value)...)
	e.buf = append(e.buf, '\t')
	return e
}

func (e *Event) Hex(label string, value []byte) *Event {
	if !e.enabled {
		return e
	}
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendHexBytes(e.buf, value)
	e.buf = append(e.buf, '\t')
	return e
}

func (e *Event) Sprintf(label, format string, a ...interface{}) *Event {
	if !e.enabled {
		return e
	}
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, escape(fmt.Sprintf(format, a...))...)
	e.buf = append(e.buf, '\t')
	return e
}

func (e *Event) Bool(label string, value bool) *Event {
	if !e.enabled {
		return e
	}
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendBool(e.buf, value)
	e.buf = append(e.buf, '\t')
	return e
}

func (e *Event) Int64(label string, value int64) *Event {
	if !e.enabled {
		return e
	}
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendInt(e.buf, value, 10)
	e.buf = append(e.buf, '\t')
	return e
}

func (e *Event) Uint64(label string, value uint64) *Event {
	if !e.enabled {
		return e
	}
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendUint(e.buf, value, 10)
	e.buf = append(e.buf, '\t')
	return e
}

func (e *Event) Float32(label string, value float32) *Event {
	if !e.enabled {
		return e
	}
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, strconv.FormatFloat(float64(value), 'g', -1, 32)...)
	e.buf = append(e.buf, '\t')
	return e
}

func (e *Event) Float64(label string, value float64) *Event {
	if !e.enabled {
		return e
	}
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = append(e.buf, strconv.FormatFloat(value, 'g', -1, 64)...)
	e.buf = append(e.buf, '\t')
	return e
}

func (e *Event) UTCTime(label string, value time.Time) *Event {
	if !e.enabled {
		return e
	}
	e.buf = append(e.buf, label...)
	e.buf = append(e.buf, ':')
	e.buf = appendUTCTime(e.buf, value)
	e.buf = append(e.buf, '\t')
	return e
}

func (e *Event) Log() {
	if e.enabled {
		e.buf[len(e.buf)-1] = '\n'
		_, _ = e.logger.writer.Write(e.buf)
	}
	eventPool.Put(e)
}
