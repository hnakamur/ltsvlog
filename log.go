// Package ltsvlog is a minimalist logging library for writing logs in
// LTSV (Labeled Tab-separated Value) format.
// See http://ltsv.org/ for LTSV.
// This library is designed with emphasis on more performance than flexibility.
//
// The first value is the current time with the label "time".
// The time format is RFC3339Nano UTC like 2006-01-02T15:04:05.999999999Z.
// The width of the nanoseconds are always 9. For example, the nanoseconds
// 123 is printed as 123000000.
// The second value is the log level with the label "level".
// Then labeled values passed to Debug, Info, Error follows.
//
// Package ltsv provides three log levels: Debug, Info and Error.
// The Info and Error levels are always enabled.
// You can disable the Debug level but only when you create a logger.
package ltsvlog

import (
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
)

// LV represents a label L and a value V.
type LV struct {
	L string
	V interface{}
}

// LTSVLogger is a LTSV logger.
type LTSVLogger struct {
	writer       io.Writer
	debugEnabled bool
	appendFunc   AppendFunc
	buf          []byte
	mu           sync.Mutex
}

// AppendFunc is a function type for appending a value to
// a byte buffer and returns the result buffer.
type AppendFunc func(buf []byte, v interface{}) []byte

// NewLTSVLogger creates a LTSV logger.
func NewLTSVLogger(w io.Writer, debugEnabled bool, appendFunc AppendFunc) *LTSVLogger {
	if appendFunc == nil {
		appendFunc = appendValue
	}
	return &LTSVLogger{
		writer:       w,
		debugEnabled: debugEnabled,
		appendFunc:   appendFunc,
	}
}

// DebugEnabled returns whether or not the debug level is enabled.
// You can avoid the cost of evaluation of arguments passed to Debug like:
// if logger.DebugEnabled() { logger.Debug(ltsvlog.LTV{"label1": "value1"}) }
func (l *LTSVLogger) DebugEnabled() bool {
	return l.debugEnabled
}

// Debug writes a log with the debug level if the debug level is enabled.
func (l *LTSVLogger) Debug(lv ...LV) {
	if l.debugEnabled {
		l.log("Debug", lv...)
	}
}

// Info writes a log with the info level.
func (l *LTSVLogger) Info(lv ...LV) {
	l.log("Info", lv...)
}

// Error writes a log with the error level.
func (l *LTSVLogger) Error(lv ...LV) {
	l.log("Error", lv...)
}

func (l *LTSVLogger) log(level string, lv ...LV) {
	l.mu.Lock()
	// Note: To reuse the buffer, create an empty slice pointing to
	// the previously allocated buffer.
	buf := append(l.buf[:0], "time:"...)
	now := time.Now().UTC()
	buf = appendZeroPaddedInt(buf, now.Year(), 1000)
	buf = append(buf, byte('-'))
	buf = appendZeroPaddedInt(buf, int(now.Month()), 10)
	buf = append(buf, byte('-'))
	buf = appendZeroPaddedInt(buf, now.Day(), 10)
	buf = append(buf, byte('T'))
	buf = appendZeroPaddedInt(buf, now.Hour(), 10)
	buf = append(buf, byte(':'))
	buf = appendZeroPaddedInt(buf, now.Minute(), 10)
	buf = append(buf, byte(':'))
	buf = appendZeroPaddedInt(buf, now.Second(), 10)
	buf = append(buf, byte('.'))
	buf = appendZeroPaddedInt(buf, now.Nanosecond(), 100000000)
	buf = append(buf, byte('Z'))

	buf = append(buf, "\tlevel:"...)
	buf = append(buf, []byte(level)...)
	for _, labelAndVal := range lv {
		buf = append(buf, '\t')
		buf = append(buf, []byte(labelAndVal.L)...)
		buf = append(buf, ':')
		buf = l.appendFunc(buf, labelAndVal.V)
	}
	buf = append(buf, '\n')
	_, _ = l.writer.Write(buf)
	l.buf = buf
	l.mu.Unlock()
}

func appendValue(buf []byte, v interface{}) []byte {
	// NOTE: In type switch case, case byte and case uint8 cannot coexist,
	// case rune and case uint cannot coexist.
	switch v.(type) {
	case nil:
		buf = append(buf, "<nil>"...)
	case string:
		buf = append(buf, []byte(v.(string))...)
	case int:
		buf = strconv.AppendInt(buf, int64(v.(int)), 10)
	case uint:
		buf = strconv.AppendUint(buf, uint64(v.(uint)), 10)
	case int8:
		buf = strconv.AppendInt(buf, int64(v.(int8)), 10)
	case int16:
		buf = strconv.AppendInt(buf, int64(v.(int16)), 10)
	case int32:
		buf = strconv.AppendInt(buf, int64(v.(int32)), 10)
	case int64:
		buf = strconv.AppendInt(buf, v.(int64), 10)
	case uint8:
		buf = strconv.AppendUint(buf, uint64(v.(uint8)), 10)
	case uint16:
		buf = strconv.AppendUint(buf, uint64(v.(uint16)), 10)
	case uint32:
		buf = strconv.AppendUint(buf, uint64(v.(uint32)), 10)
	case uint64:
		buf = strconv.AppendUint(buf, uint64(v.(uintptr)), 10)
	case float32:
		buf = append(buf, []byte(strconv.FormatFloat(float64(v.(float32)), 'g', -1, 32))...)
	case float64:
		buf = append(buf, []byte(strconv.FormatFloat(v.(float64), 'g', -1, 64))...)
	case bool:
		buf = strconv.AppendBool(buf, v.(bool))
	case uintptr:
		buf = strconv.AppendUint(buf, uint64(v.(uintptr)), 10)
	case []byte:
		buf = appendHexBytes(buf, v.([]byte))
	case fmt.Stringer:
		buf = append(buf, []byte(v.(fmt.Stringer).String())...)
	default:
		buf = append(buf, []byte(fmt.Sprintf("%v", v))...)
	}
	return buf
}

var digits = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

func appendHexBytes(buf []byte, v []byte) []byte {
	buf = append(buf, "0x"...)
	for _, b := range v {
		buf = append(buf, digits[b/16])
		buf = append(buf, digits[b%16])
	}
	return buf
}

func appendZeroPaddedInt(buf []byte, v, p int) []byte {
	for ; p > 1; p /= 10 {
		q := v / p
		buf = append(buf, digits[q])
		v -= q * p
	}
	return append(buf, digits[v])
}
