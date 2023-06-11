package ltsvlog

import (
	"io"
	"sync"
)

var _ io.Writer = (*SwitchableWriter)(nil)

type SwitchableWriter struct {
	w  io.Writer
	mu sync.Mutex
}

func NewSwitchableWriter(w io.Writer) *SwitchableWriter {
	return &SwitchableWriter{w: w}
}

func (sw *SwitchableWriter) Write(p []byte) (n int, err error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	return sw.w.Write(p)
}

func (sw *SwitchableWriter) Switch(w io.Writer) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	sw.w = w
}
