package ltsvlog

import (
	"bytes"
	"runtime"
)

// Stack formats a stack trace of the calling goroutine into buf
// in one line format which suitable for LTSV logs.
// It returns the string converted from buf. If buf is small,
// only the partial stack trace is returned. If you pass nil
// as buf, a byte buffer of size 8192 is used internally.
func Stack(buf []byte) string {
	if buf == nil {
		buf = make([]byte, 1024)
	}
	n := runtime.Stack(buf, false)
	buf = buf[:n]

	end := bytes.IndexByte(buf, '\n')
	if end == -1 {
		return ""
	}
	p := append(buf[:0], buf[:end]...)

	// Skip the first stack since it is this function.
	start := indexByteFrom(buf, '\n', end+1)
	if start == -1 {
		return string(p)
	}
	start = indexByteFrom(buf, '\n', end+1+start+1)
	if start == -1 {
		return string(p)
	}
	start++

	for start < n {
		p = append(p, " ["...)
		end = indexByteFrom(buf, '\n', start)
		if end == -1 {
			p = append(p, buf[start:n]...)
			p = append(p, "..."...)
			break
		}
		p = append(p, buf[start:end]...)
		p = append(p, ' ')
		start = end + 2
		if start >= n {
			p = append(p, "..."...)
			break
		}
		end = indexByteFrom(buf, '\n', start)
		if end == -1 {
			p = append(p, buf[start:n]...)
			p = append(p, "..."...)
			break
		}
		p = append(p, buf[start:end]...)
		p = append(p, ']')
		start = end + 1
	}
	return string(p)
}

func indexByteFrom(buf []byte, b byte, offset int) int {
	i := bytes.IndexByte(buf[offset:], b)
	if i == -1 {
		return -1
	}
	return i + offset
}
