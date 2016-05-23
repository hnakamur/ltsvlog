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
// If you pass nil as buf, a byte buffer of size 8192 is used internally.
func Stack(buf []byte) string {
	if buf == nil {
		buf = make([]byte, 8192)
	}
	n := runtime.Stack(buf, false)
	buf = buf[:n]

	// NOTE: We reuse the same buffer here.
	p := buf[:0]

	i := bytes.IndexByte(buf, '\n')
	if i == -1 {
		goto buffer_too_small
	}

	// NOTE: Skip the first stack since it is this function.
	if i+1 > len(buf) {
		goto buffer_too_small
	}
	buf = buf[i+1:]

	i = bytes.IndexByte(buf, '\n')
	if i == -1 || i+1 > len(buf) {
		goto buffer_too_small
	}
	buf = buf[i+1:]

	i = bytes.IndexByte(buf, '\n')
	if i == -1 || i+1 > len(buf) {
		goto buffer_too_small
	}
	buf = buf[i+1:]

	for len(buf) > 0 {
		p = append(p, '[')
		i = bytes.IndexByte(buf, '\n')
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
