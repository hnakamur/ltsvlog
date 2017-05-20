package ltsvlog

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var lvsPool = &sync.Pool{
	New: func() interface{} {
		return &LVs{
			buf: make([]byte, 8192),
		}
	},
}

type LVs struct {
	logger  *LTSVLogger
	enabled bool
	buf     []byte
}

func (v *LVs) String(label string, value string) *LVs {
	if !v.enabled {
		return v
	}
	if len(v.buf) > 0 {
		v.buf = append(v.buf, '\t')
	}
	v.buf = append(v.buf, label...)
	v.buf = append(v.buf, ':')
	v.buf = append(v.buf, escape(value)...)
	return v
}

func (v *LVs) Hex(label string, value []byte) *LVs {
	if !v.enabled {
		return v
	}
	if len(v.buf) > 0 {
		v.buf = append(v.buf, '\t')
	}
	v.buf = append(v.buf, label...)
	v.buf = append(v.buf, ':')
	v.buf = appendHexBytes(v.buf, value)
	return v
}

func (v *LVs) Sprintf(label, format string, a ...interface{}) *LVs {
	if !v.enabled {
		return v
	}
	if len(v.buf) > 0 {
		v.buf = append(v.buf, '\t')
	}
	v.buf = append(v.buf, label...)
	v.buf = append(v.buf, ':')
	v.buf = append(v.buf, escape(fmt.Sprintf(format, a...))...)
	return v
}

func (v *LVs) Bool(label string, value bool) *LVs {
	if !v.enabled {
		return v
	}
	if len(v.buf) > 0 {
		v.buf = append(v.buf, '\t')
	}
	v.buf = append(v.buf, label...)
	v.buf = append(v.buf, ':')
	v.buf = strconv.AppendBool(v.buf, value)
	return v
}

func (v *LVs) Int64(label string, value int64) *LVs {
	if !v.enabled {
		return v
	}
	if len(v.buf) > 0 {
		v.buf = append(v.buf, '\t')
	}
	v.buf = append(v.buf, label...)
	v.buf = append(v.buf, ':')
	v.buf = strconv.AppendInt(v.buf, value, 10)
	return v
}

func (v *LVs) Uint64(label string, value uint64) *LVs {
	if !v.enabled {
		return v
	}
	if len(v.buf) > 0 {
		v.buf = append(v.buf, '\t')
	}
	v.buf = append(v.buf, label...)
	v.buf = append(v.buf, ':')
	v.buf = strconv.AppendUint(v.buf, value, 10)
	return v
}

func (v *LVs) Float32(label string, value float32) *LVs {
	if !v.enabled {
		return v
	}
	if len(v.buf) > 0 {
		v.buf = append(v.buf, '\t')
	}
	v.buf = append(v.buf, label...)
	v.buf = append(v.buf, ':')
	v.buf = append(v.buf, strconv.FormatFloat(float64(value), 'g', -1, 32)...)
	return v
}

func (v *LVs) Float64(label string, value float64) *LVs {
	if !v.enabled {
		return v
	}
	if len(v.buf) > 0 {
		v.buf = append(v.buf, '\t')
	}
	v.buf = append(v.buf, label...)
	v.buf = append(v.buf, ':')
	v.buf = append(v.buf, strconv.FormatFloat(value, 'g', -1, 64)...)
	return v
}

func (v *LVs) UTCTime(label string, value time.Time) *LVs {
	if !v.enabled {
		return v
	}
	if len(v.buf) > 0 {
		v.buf = append(v.buf, '\t')
	}
	v.buf = append(v.buf, label...)
	v.buf = append(v.buf, ':')
	v.buf = appendUTCTime(v.buf, value)
	return v
}

func (v *LVs) Log() {
	if v.enabled {
		v.buf = append(v.buf, '\n')
		_, _ = v.logger.writer.Write(v.buf)
	}
	lvsPool.Put(v)
}
