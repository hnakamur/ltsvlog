// Package ltsvlog is a minimalist logging library for writing logs in
// LTSV (Labeled Tab-separated Value) format.
// See http://ltsv.org/ for LTSV.
//
// This logging library has three log levels: Debug, Info and Error.
// The Info and Error levels are always enabled.
// You can disable the Debug level but only when you create a logger.
//
// Each log record is printed as one line. A line has multiple fields
// separated by a tab character. Each field has a label and a value
// separated by a colon ':' character.
// So you must not contain a new line or a tab character in your labels
// and values. You must not contain a colon character in your labels.
package ltsvlog

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// LV represents a label L and a value V.
type LV struct {
	L string
	V interface{}
}

// LogWriter is a LTSV logger interface
type LogWriter interface {
	DebugEnabled() bool
	Debug(lv ...LV)
	Info(lv ...LV)
	Error(lv ...LV)
	ErrorWithStack(lv ...LV)
}

// LTSVLogger is a LTSV logger.
type LTSVLogger struct {
	writer                  io.Writer
	debugEnabled            bool
	timeLabel               string
	levelLabel              string
	appendPrefixFunc        AppendPrefixFunc
	appendPrefixFuncChanged bool
	appendValueFunc         AppendValueFunc
	buf                     []byte
	stackBuf                []byte
	mu                      sync.Mutex
}

// Option is the function type to set an option of LTSVLogger
type Option func(l *LTSVLogger)

// StackBufSize returns the option function to set the stack buffer size.
func StackBufSize(size int) Option {
	return func(l *LTSVLogger) {
		l.stackBuf = make([]byte, size)
	}
}

// SetTimeLabel returns the option function to set the time label.
// If the label is empty, logger does not print time values.
func SetTimeLabel(label string) Option {
	return func(l *LTSVLogger) {
		l.timeLabel = label
	}
}

// SetLevelLabel returns the option function to set the level label.
// If the label is empty, logger does not print level values.
func SetLevelLabel(label string) Option {
	return func(l *LTSVLogger) {
		l.levelLabel = label
	}
}

// SetAppendPrefix returns the option function to set the function
// to append the log prefix. In the default setting, the prefix consists
// of the time and the level.
func SetAppendPrefix(f AppendPrefixFunc) Option {
	return func(l *LTSVLogger) {
		l.appendPrefixFunc = f
		l.appendPrefixFuncChanged = true
	}
}

// SetAppendValue returns the option function to set the function
// to append a value.
func SetAppendValue(f AppendValueFunc) Option {
	return func(l *LTSVLogger) {
		l.appendValueFunc = f
	}
}

// AppendPrefixFunc is a function type for appending a prefix
// for a log record to a byte buffer and returns the result buffer.
type AppendPrefixFunc func(buf []byte, level string) []byte

// AppendValueFunc is a function type for appending a value to
// a byte buffer and returns the result buffer.
type AppendValueFunc func(buf []byte, v interface{}) []byte

const (
	defaultTimeLabel  = "time"
	defaultLevelLabel = "level"
)

var defaultAppendPrefixFunc = appendPrefixFunc(defaultTimeLabel, defaultLevelLabel)

// NewLTSVLogger creates a LTSV logger with the default time and value format.
// Shorthand for NewLTSVLoggerCustomFormat(w, debugEnabled, 8192, nil, nil).
//
// The folloing two values are prepended to each log line.
//
// The first value is the current time with the label "time".
// The time format is RFC3339Nano UTC like 2006-01-02T15:04:05.00000000Z.
// The width of the nanoseconds are always 9. For example, the nanoseconds
// 123 is printed as 123000000.
// The second value is the log level with the label "level".
func NewLTSVLogger(w io.Writer, debugEnabled bool, options ...Option) *LTSVLogger {
	l := &LTSVLogger{
		writer:           w,
		debugEnabled:     debugEnabled,
		timeLabel:        defaultTimeLabel,
		levelLabel:       defaultLevelLabel,
		appendPrefixFunc: defaultAppendPrefixFunc,
		appendValueFunc:  appendValue,
		stackBuf:         make([]byte, 8192),
	}
	for _, o := range options {
		o(l)
	}
	if !l.appendPrefixFuncChanged &&
		(l.timeLabel != defaultTimeLabel || l.levelLabel != defaultLevelLabel) {
		l.appendPrefixFunc = appendPrefixFunc(l.timeLabel, l.levelLabel)
	}
	return l
}

// Deprecated. Use NewLTSVLogger with options instead.
//
// NewLTSVLoggerCustomFormat creates a LTSV logger with the buffer size for
// filling stack traces and user-supplied functions for appending a log
// record prefix and appending a log value.
func NewLTSVLoggerCustomFormat(w io.Writer, debugEnabled bool, stackBufSize int, appendPrefixFunc AppendPrefixFunc, appendValueFunc AppendValueFunc) *LTSVLogger {
	if appendPrefixFunc == nil {
		appendPrefixFunc = appendPrefix
	}
	if appendValueFunc == nil {
		appendValueFunc = appendValue
	}
	return &LTSVLogger{
		writer:           w,
		debugEnabled:     debugEnabled,
		appendPrefixFunc: appendPrefixFunc,
		appendValueFunc:  appendValueFunc,
		stackBuf:         make([]byte, stackBufSize),
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
		l.mu.Lock()
		l.log("Debug", lv...)
		l.mu.Unlock()
	}
}

// Info writes a log with the info level.
func (l *LTSVLogger) Info(lv ...LV) {
	l.mu.Lock()
	l.log("Info", lv...)
	l.mu.Unlock()
}

// Error writes a log with the error level.
func (l *LTSVLogger) Error(lv ...LV) {
	l.mu.Lock()
	l.log("Error", lv...)
	l.mu.Unlock()
}

// ErrorWithStack writes a log and a stack with the error level.
func (l *LTSVLogger) ErrorWithStack(lv ...LV) {
	l.mu.Lock()
	args := lv
	args = append(args, LV{"stack", stack(2, l.stackBuf)})
	l.log("Error", args...)
	l.mu.Unlock()
}

func (l *LTSVLogger) log(level string, lv ...LV) {
	// Note: To reuse the buffer, create an empty slice pointing to
	// the previously allocated buffer.
	buf := l.appendPrefixFunc(l.buf[:0], level)
	for i, labelAndVal := range lv {
		if i > 0 {
			buf = append(buf, '\t')
		}
		buf = append(buf, []byte(labelAndVal.L)...)
		buf = append(buf, ':')
		buf = l.appendValueFunc(buf, labelAndVal.V)
	}
	buf = append(buf, '\n')
	_, _ = l.writer.Write(buf)
	l.buf = buf
}

func appendPrefixFunc(timeLabel, levelLabel string) AppendPrefixFunc {
	if timeLabel != "" && levelLabel != "" {
		timeLabelBytes := []byte(timeLabel + ":")
		levelLabelBytes := []byte("\t" + levelLabel + ":")
		return func(buf []byte, level string) []byte {
			buf = append(buf, timeLabelBytes...)
			now := time.Now().UTC()
			buf = appendTime(buf, now)
			buf = append(buf, levelLabelBytes...)
			buf = append(buf, []byte(level)...)
			buf = append(buf, '\t')
			return buf
		}
	} else if timeLabel != "" && levelLabel == "" {
		timeLabelBytes := []byte(timeLabel + ":")
		return func(buf []byte, level string) []byte {
			buf = append(buf, timeLabelBytes...)
			now := time.Now().UTC()
			buf = appendTime(buf, now)
			buf = append(buf, '\t')
			return buf
		}
	} else if timeLabel == "" && levelLabel != "" {
		levelLabelBytes := []byte(levelLabel + ":")
		return func(buf []byte, level string) []byte {
			buf = append(buf, levelLabelBytes...)
			buf = append(buf, []byte(level)...)
			buf = append(buf, '\t')
			return buf
		}
	} else {
		return func(buf []byte, level string) []byte {
			return buf
		}
	}
}

func appendPrefix(buf []byte, level string) []byte {
	buf = append(buf, "time:"...)
	now := time.Now().UTC()
	buf = appendTime(buf, now)
	buf = append(buf, "\tlevel:"...)
	buf = append(buf, []byte(level)...)
	buf = append(buf, '\t')
	return buf
}

func appendTime(buf []byte, t time.Time) []byte {
	buf = appendZeroPaddedInt(buf, t.Year(), 1000)
	buf = append(buf, byte('-'))
	buf = appendZeroPaddedInt(buf, int(t.Month()), 10)
	buf = append(buf, byte('-'))
	buf = appendZeroPaddedInt(buf, t.Day(), 10)
	buf = append(buf, byte('T'))
	buf = appendZeroPaddedInt(buf, t.Hour(), 10)
	buf = append(buf, byte(':'))
	buf = appendZeroPaddedInt(buf, t.Minute(), 10)
	buf = append(buf, byte(':'))
	buf = appendZeroPaddedInt(buf, t.Second(), 10)
	buf = append(buf, byte('.'))
	buf = appendZeroPaddedInt(buf, t.Nanosecond(), 100000000)
	return append(buf, byte('Z'))
}

var escaper = strings.NewReplacer("\t", "\\t", "\n", "\\n")

func escape(s string) string {
	return escaper.Replace(s)
}

func appendValue(buf []byte, v interface{}) []byte {
	// NOTE: In type switch case, case byte and case uint8 cannot coexist,
	// case rune and case uint cannot coexist.
	switch v.(type) {
	case nil:
		buf = append(buf, "<nil>"...)
	case string:
		buf = append(buf, []byte(escape(v.(string)))...)
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
		buf = strconv.AppendUint(buf, uint64(v.(uint64)), 10)
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
		buf = append(buf, []byte(escape(v.(fmt.Stringer).String()))...)
	default:
		buf = append(buf, []byte(escape(fmt.Sprintf("%+v", v)))...)
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

// Logger is the global logger.
// You can change this logger like
// ltsvlog.Logger = ltsvlog.NewLTSVLogger(os.Stdout, false)
// You can change the global logger safely only before writing
// to the logger. Chaging the logger while writing may cause
// the unexpected behavior.
var Logger = NewLTSVLogger(os.Stdout, true)

// Discard discards any logging outputs.
type Discard struct{}

// DebugEnabled always return false
func (*Discard) DebugEnabled() bool { return false }

// Debug prints nothing
func (*Discard) Debug(lv ...LV) {}

// Info prints nothing
func (*Discard) Info(lv ...LV) {}

// Error prints nothing
func (*Discard) Error(lv ...LV) {}

// ErrorWithStack prints nothing
func (*Discard) ErrorWithStack(lv ...LV) {}
