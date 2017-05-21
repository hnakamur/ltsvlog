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
		ltsvlog.Logger.Debug().String("msg", "This is a debug message").
			String("str", "foo").Int("int", 234).Log()
	}

	ltsvlog.Logger.Info().Sprintf("float1", "%3.2f", 3.14).Log()

	err := a()
	if err != nil {
		ltsvlog.Logger.Err(err)
	}
}

func a() error {
	return b()
}

func b() error {
	return ltsvlog.Err(errors.New("some error")).String("key1", "value1").Stack("")
}
```

An example output:

```
time:2017-05-21T05:19:11.256860Z	level:Debug	msg:This is a debug message	str:foo	int:234
time:2017-05-21T05:19:11.256887Z	level:Info	float1:3.14
time:2017-05-21T05:19:11.256926Z	level:Error	err:some error	key1:value1	stack:[main.b(0x0, 0x0) /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/main.go:28 +0xc8],[main.a(0xc42001a240, 0x4c5d6e) /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/main.go:24 +0x22],[main.main() /home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/main.go:17 +0x11e]
```

Since these log lines ar long, please scroll horizontally to the right to see all the output.

## Benchmark result
[hnakamur/go-log-benchmarks](https://github.com/hnakamur/go-log-benchmarks)

## License
MIT
