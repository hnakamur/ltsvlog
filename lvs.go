package ltsvlog

import "sync"

var lvsPool = &sync.Pool{
	New: func() interface{} {
		return &LVs{
			lvs: make([]LV, 0, 8),
		}
	},
}

// LVs holds a log level and labeled values.
type LVs struct {
	level string
	lvs   []LV
}

// Info creates a LVs with the Info level.
func Info() *LVs {
	lvs := lvsPool.Get().(*LVs)
	lvs.level = "Info"
	lvs.lvs = lvs.lvs[:0]
	return lvs
}

// Info creates a LVs with the Debug level.
func Debug() *LVs {
	lvs := lvsPool.Get().(*LVs)
	lvs.level = "Debug"
	lvs.lvs = lvs.lvs[:0]
	return lvs
}

// LV appends a label value pair to LVs.
func (v *LVs) LV(key string, value interface{}) *LVs {
	v.lvs = append(v.lvs, LV{key, value})
	return v
}
