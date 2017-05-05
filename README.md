ltsvlog [![Build Status](https://travis-ci.org/hnakamur/ltsvlog.png)](https://travis-ci.org/hnakamur/ltsvlog) [![Go Report Card](https://goreportcard.com/badge/github.com/hnakamur/ltsvlog)](https://goreportcard.com/report/github.com/hnakamur/ltsvlog) [![GoDoc](https://godoc.org/github.com/hnakamur/ltsvlog?status.svg)](https://godoc.org/github.com/hnakamur/ltsvlog) [![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hyperium/hyper/master/LICENSE)
=======

ltsvlog is a minimalist [LTSV; Labeled Tab-separated Values](http://ltsv.org/) logging library in Go.
See https://godoc.org/github.com/hnakamur/ltsvlog for the API document.

I wrote a blog article about this library in Japanese: [GoでLTSV形式でログ出力するライブラリを書いた · hnakamur's blog at github](http://hnakamur.github.io/blog/2016/06/13/wrote_go_ltsvlog_library/).

## An example code and output

An example code:

```
package main

import (
	"errors"

	"github.com/hnakamur/ltsvlog"
)

func main() {
	if ltsvlog.Logger.DebugEnabled() {
		ltsvlog.Logger.Debug(ltsvlog.LV{"msg", "This is a debug message"},
			ltsvlog.LV{"key", "key1"}, ltsvlog.LV{"intValue", 234})
	}
	ltsvlog.Logger.Info(ltsvlog.LV{"msg", "hello, world"}, ltsvlog.LV{"key", "key1"},
		ltsvlog.LV{"value", "value1"})
	a()
	ltsvlog.Logger.Info(ltsvlog.LV{"msg", "goodbye, world"}, ltsvlog.LV{"foo", "bar"},
		ltsvlog.LV{"nilValue", nil}, ltsvlog.LV{"bytes", []byte("a/b")})
}

func a() {
	b()
}

func b() {
	err := errors.New("demo error")
	if err != nil {
		ltsvlog.Logger.ErrorWithStack(ltsvlog.LV{"err", err})
	}
}
```

An example output:

```
$ go run cmd/example/main.go
time:2016-05-30T02:21:28.135713584Z     level:Debug     msg:This is a debug message     key:key1        intValue:234
time:2016-05-30T02:21:28.135744631Z     level:Info      msg:hello, world        key:key1        value:value1
time:2016-05-30T02:21:28.135772957Z     level:Error     err:demo error  stack:[main.b() /home/hnakamur/gocode/src/github.com/hnakamur/ltsvlog/cmd/example/main.go:28 +0x1ba],[main.a() /home/hnakamur/gocode/src/github.com/hnakamur/ltsvlog/cmd/example/main.go:22 +0x14],[main.main() /home/hnakamur/gocode/src/github.com/hnakamur/ltsvlog/cmd/example/main.go:16 +0x56c]
time:2016-05-30T02:21:28.135804911Z     level:Info      msg:goodbye, world      foo:bar nilValue:<nil>  bytes:0x612f62
```

Since these log lines ar long, please scroll horizontally to the right to see all the output.

## Benchmark result

```
$ go test -count=10 -bench . -benchmem
BenchmarkLTSVLog-2               2000000               700 ns/op             238 B/op          3 allocs/op
BenchmarkLTSVLog-2               2000000               705 ns/op             238 B/op          3 allocs/op
BenchmarkLTSVLog-2               2000000               702 ns/op             238 B/op          3 allocs/op
BenchmarkLTSVLog-2               2000000               703 ns/op             238 B/op          3 allocs/op
BenchmarkLTSVLog-2               2000000               704 ns/op             238 B/op          3 allocs/op
BenchmarkLTSVLog-2               2000000               706 ns/op             238 B/op          3 allocs/op
BenchmarkLTSVLog-2               2000000               714 ns/op             238 B/op          3 allocs/op
BenchmarkLTSVLog-2               2000000               705 ns/op             238 B/op          3 allocs/op
BenchmarkLTSVLog-2               2000000               703 ns/op             238 B/op          3 allocs/op
BenchmarkLTSVLog-2               2000000               703 ns/op             238 B/op          3 allocs/op
BenchmarkStandardLog-2           2000000               787 ns/op             274 B/op          3 allocs/op
BenchmarkStandardLog-2           2000000               790 ns/op             274 B/op          3 allocs/op
BenchmarkStandardLog-2           2000000               790 ns/op             274 B/op          3 allocs/op
BenchmarkStandardLog-2           2000000               790 ns/op             274 B/op          3 allocs/op
BenchmarkStandardLog-2           2000000               794 ns/op             274 B/op          3 allocs/op
BenchmarkStandardLog-2           2000000               790 ns/op             274 B/op          3 allocs/op
BenchmarkStandardLog-2           2000000               790 ns/op             274 B/op          3 allocs/op
BenchmarkStandardLog-2           2000000               793 ns/op             274 B/op          3 allocs/op
BenchmarkStandardLog-2           2000000               792 ns/op             274 B/op          3 allocs/op
BenchmarkStandardLog-2           2000000               789 ns/op             274 B/op          3 allocs/op
PASS
ok      github.com/hnakamur/ltsvlog     45.224s
```

## License
MIT
