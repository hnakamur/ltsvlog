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
	buf = appendTime(buf, now)
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

func appendTime(buf []byte, t time.Time) []byte {
	buf = appendDigits(buf, t.Year(), 4)
	buf = append(buf, byte('-'))
	buf = appendDigits(buf, int(t.Month()), 2)
	buf = append(buf, byte('-'))
	buf = appendDigits(buf, t.Day(), 2)
	buf = append(buf, byte('T'))
	buf = appendDigits(buf, t.Hour(), 2)
	buf = append(buf, byte(':'))
	buf = appendDigits(buf, t.Minute(), 2)
	buf = append(buf, byte(':'))
	buf = appendDigits(buf, t.Second(), 2)
	buf = append(buf, byte('.'))
	buf = appendDigits(buf, t.Nanosecond(), 9)
	return append(buf, byte('Z'))
}

func appendDigits(buf []byte, value, width int) []byte {
	width--
	p := pow10(width)
	for ; width > 0; width-- {
		q := value / p
		buf = append(buf, digits[q])
		value -= q * p
		p /= 10
	}
	return append(buf, digits[value])
}

var digits = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

func pow10(e int) int {
	p := 1
	for ; e > 0; e-- {
		p *= 10
	}
	return p
}
