package ltsvlog

import "sync"

var lvsPool = &sync.Pool{
	New: func() interface{} {
		return &LVs{
			buf: make([]byte, 8192),
		}
	},
}

type LVs struct {
	logger *LTSVLogger
	buf    []byte
}

func (v *LVs) LV(key string, value interface{}) *LVs {
	v.buf = append(v.buf, '\t')
	v.buf = append(v.buf, key...)
	v.buf = append(v.buf, ':')
	v.buf = v.logger.appendValueFunc(v.buf, value)
	return v
}

func (v *LVs) Debug() {
	l := v.logger
	l.mu.Lock()
	l.rawLog("Debug", v.buf)
	l.mu.Unlock()
	lvsPool.Put(v)
}

func (v *LVs) Info() {
	l := v.logger
	l.mu.Lock()
	l.rawLog("Info", v.buf)
	l.mu.Unlock()
	lvsPool.Put(v)
}
