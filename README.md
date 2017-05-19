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
time:2017-05-05T14:55:30.326842Z        level:Debug     msg:This is a debug message     key:key1       intValue:234
time:2017-05-05T14:55:30.326858Z        level:Info      msg:hello, world        key:key1      value:value1
time:2017-05-05T14:55:30.326905Z        level:Error     err:demo error  stack:[main.b() /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/cmd/example/main.go:28 +0xd6],[main.a() /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/cmd/example/main.go:22 +0x20],[main.main() /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/cmd/example/main.go:16 +0x1fb]
time:2017-05-05T14:55:30.326937Z        level:Info      msg:goodbye, world      foo:bar nilValue:<nil> bytes:0x612f62
```

Since these log lines ar long, please scroll horizontally to the right to see all the output.

## Benchmark result
[hnakamur/go-log-benchmarks](https://github.com/hnakamur/go-log-benchmarks)

## License
MIT
