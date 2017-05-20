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

// Event is a temporary object for building a log record of
// Debug or Info level.
type Event struct {
	logger  *LTSVLogger
	enabled bool
	buf     []byte
}

// String appends a labeled string value to Event.
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

// Hex appends a labeled hex value to Event.
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

// Sprintf appends a labeled formatted string value to Event.
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

// Bool appends a labeled bool value to Event.
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

// Int64 appends a labeled int64 value to Event.
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

// Uint64 appends a labeled uint64 value to Event.
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

// Float32 appends a labeled float32 value to Event.
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

// Float64 appends a labeled float64 value to Event.
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

// UTCTime appends a labeled UTC time value to Event.
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

// Log write this event if the logger which created this event is enabled.
func (e *Event) Log() {
	if e.enabled {
		e.buf[len(e.buf)-1] = '\n'
		_, _ = e.logger.writer.Write(e.buf)
	}
	eventPool.Put(e)
}
