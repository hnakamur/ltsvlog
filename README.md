ltsvlog
=======

ltsvlog is a minimalist [LTSV; Labeled Tab-separated Values](http://ltsv.org/) logging library in Go.
See https://godoc.org/github.com/hnakamur/ltsvlog for the API document.

## An example code and output

An example code:

```
package main

import (
	"errors"
	"os"

	"github.com/hnakamur/ltsvlog"
)

var logger *ltsvlog.LTSVLogger

func main() {
	logger = ltsvlog.NewLTSVLogger(os.Stdout, true)
	if logger.DebugEnabled() {
		logger.Debug(ltsvlog.LV{"msg", "This is a debug message"},
			ltsvlog.LV{"key", "key1"}, ltsvlog.LV{"intValue", 234})
	}
	logger.Info(ltsvlog.LV{"msg", "hello, world"}, ltsvlog.LV{"key", "key1"},
		ltsvlog.LV{"value", "value1"})
	a()
	logger.Info(ltsvlog.LV{"msg", "goodbye, world"}, ltsvlog.LV{"foo", "bar"},
		ltsvlog.LV{"nilValue", nil}, ltsvlog.LV{"bytes", []byte("a/b")})
}

func a() {
	b()
}

func b() {
	err := errors.New("demo error")
	if err != nil {
		logger.ErrorWithStack(ltsvlog.LV{"err", err})
	}
}
```

An example output:

```
$ go run cmd/example/main.go
time:2016-05-27T06:51:10.977296010Z     level:Debug     msg:This is a debug message     key:key1        intValue:234
time:2016-05-27T06:51:10.977322761Z     level:Info      msg:hello, world        key:key1        value:value1
time:2016-05-27T06:51:10.977357972Z     level:Error     err:demo error  stack:[main.b() /home/hnakamur/gocode/src/github.com/hnakamur/ltsvlog/cmd/example/main.go:32 +0x1ba],[main.a() /home/hnakamur/gocode/src/github.com/hnakamur/ltsvlog/cmd/example/main.go:26 +0x14],[main.main() /home/hnakamur/gocode/src/github.com/hnakamur/ltsvlog/cmd/example/main.go:20 +0x762]
time:2016-05-27T06:51:10.977384889Z     level:Info      msg:goodbye, world      foo:bar nilValue:<nil>  bytes:0x612f62
```

Since these log lines ar long, please scroll horizontally to the right to see all the output.

## Benchmark result

```
$ go test -bench . -benchmem
testing: warning: no tests to run
PASS
BenchmarkLTSVLog-2       1000000              1266 ns/op             245 B/op          3 allocs/op
BenchmarkStandardLog-2   1000000              1212 ns/op             235 B/op          3 allocs/op
ok      github.com/hnakamur/ltsvlog     2.542s
```

## License
MIT
