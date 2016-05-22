// Package ltsvlog is a minimalist logging library for writing logs in
// LTSV (Labeled Tab-separated Value) format.
// See http://ltsv.org/ for LTSV.
// This library is designed with emphasis on more performance than flexibility.
//
// The first value is the current time with the label "time".
// The time format is RFC3339Nano UTC like 2006-01-02T15:04:05.999999999Z.
// The second value is the log level with the label "level".
// Then labeled values passed to Debug, Info, Error follows.
//
// Package ltsv provides three log levels: Debug, Info and Error.
// The Info and Error levels are always enabled.
// You can disable the Debug level but only when you create a logger.
package ltsvlog

import (
	"io"
	"sync"
	"time"
)

// LV represents a label L and a value V.
type LV struct {
	L string
	V string
}

// LTSVLogger is a LTSV logger.
type LTSVLogger struct {
	writer       io.Writer
	debugEnabled bool
	buf          []byte
	mu           sync.Mutex
}

// NewLTSVLogger creates a LTSV logger.
func NewLTSVLogger(w io.Writer, debugEnabled bool) *LTSVLogger {
	return &LTSVLogger{
		writer:       w,
		debugEnabled: debugEnabled,
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
	buf := append(l.buf[:0], []byte("time:")...)
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

	buf = append(buf, []byte("\tlevel:")...)
	buf = append(buf, []byte(level)...)
	for _, labelAndVal := range lv {
		buf = append(buf, byte('\t'))
		buf = append(buf, []byte(labelAndVal.L)...)
		buf = append(buf, byte(':'))
		buf = append(buf, []byte(labelAndVal.V)...)
	}
	buf = append(buf, byte('\n'))
	_, _ = l.writer.Write(buf)
	l.buf = buf
	l.mu.Unlock()
}

var digits = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

func appendZeroPaddedInt(buf []byte, v, p int) []byte {
	for ; p > 1; p /= 10 {
		q := v / p
		buf = append(buf, digits[q])
		v -= q * p
	}
	return append(buf, digits[v])
}
