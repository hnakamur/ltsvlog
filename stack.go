package ltsvlog

import (
	"bytes"
	"runtime"
)

// Stack formats a stack trace of the calling goroutine into buf
// in one line format which suitable for LTSV logs.
// It returns the string converted from buf. If buf is small,
// only the partial stack trace followed by "..." is returned.
// It is the caller's responsibility to pass a large enough buf.
func stack(skip int, buf []byte) string {
	n := runtime.Stack(buf, false)
	buf = buf[:n]

	// NOTE: We reuse the same buffer here.
	p := buf[:0]

	for j := 0; j < 1+2*skip; j++ {
		i := bytes.IndexByte(buf, '\n')
		if i == -1 || i+1 > len(buf) {
			goto buffer_too_small
		}
		buf = buf[i+1:]
	}

	for len(buf) > 0 {
		p = append(p, '[')
		i := bytes.IndexByte(buf, '\n')
		if i == -1 {
			goto buffer_too_small
		}
		p = append(p, buf[:i]...)
		p = append(p, ' ')
		if i+2 > len(buf) {
			goto buffer_too_small
		}
		buf = buf[i+2:]
		i = bytes.IndexByte(buf, '\n')
		if i == -1 {
			goto buffer_too_small
		}
		p = append(p, buf[:i]...)
		p = append(p, ']')
		if i+1 > len(buf) {
			goto buffer_too_small
		}
		buf = buf[i+1:]
		if len(buf) > 0 {
			p = append(p, ',')
		}
	}
	return string(p)

buffer_too_small:
	p = append(p, buf...)
	p = append(p, "..."...)
	return string(p)
}
