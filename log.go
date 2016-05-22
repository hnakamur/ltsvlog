package ltsvlog

import (
	"io"
	"sync"
	"time"
)

type LV struct {
	L string
	V string
}

type LTSVLogger struct {
	writer       io.Writer
	buf          []byte
	debugEnabled bool
	mu           sync.Mutex
}

func NewLTSVLogger(w io.Writer, debugEnabled bool) *LTSVLogger {
	return &LTSVLogger{
		writer:       w,
		debugEnabled: debugEnabled,
	}
}

func (l *LTSVLogger) DebugEnabled() bool {
	return l.debugEnabled
}

func (l *LTSVLogger) Debug(lv ...LV) {
	if l.debugEnabled {
		l.log("Debug", lv...)
	}
}

func (l *LTSVLogger) Info(lv ...LV) {
	l.log("Info", lv...)
}

func (l *LTSVLogger) Error(lv ...LV) {
	l.log("Error", lv...)
}

func (l *LTSVLogger) log(level string, lv ...LV) {
	now := time.Now().UTC()
	l.mu.Lock()

	buf := append(l.buf[:0], []byte("time:")...)
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

func appendZeroPaddedInt(buf []byte, value, p int) []byte {
	for ; p > 1; p /= 10 {
		q := value / p
		buf = append(buf, digits[q])
		value -= q * p
	}
	return append(buf, digits[value])
}

var digits = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
