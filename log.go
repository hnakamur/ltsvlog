package ltsvlog

import (
	"sync"
	"time"
)

type LV struct {
	L string
	V string
}

type Logger interface {
	Log(v interface{})
}

type LTSVLogger struct {
	logger       Logger
	buf          []byte
	debugEnabled bool
	mu           sync.Mutex
}

func NewLTSVLogger(logger Logger, debugEnabled bool) *LTSVLogger {
	return &LTSVLogger{
		debugEnabled: debugEnabled,
		logger:       logger,
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
	l.mu.Lock()
	now := time.Now().Format(time.RFC3339Nano)
	buf := append(l.buf[:0], []byte("time:")...)
	buf = append(buf, []byte(now)...)
	buf = append(buf, []byte("\tlevel:")...)
	buf = append(buf, []byte(level)...)
	for _, labelAndVal := range lv {
		buf = append(buf, byte('\t'))
		buf = append(buf, []byte(labelAndVal.L)...)
		buf = append(buf, byte(':'))
		buf = append(buf, []byte(labelAndVal.V)...)
	}
	l.logger.Log(string(buf))
	l.buf = buf
	l.mu.Unlock()
}
