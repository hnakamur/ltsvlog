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

Note these benchmarks print roughly same outputs, but not the exactly same outputs.

Especially BenchmarkZapLTSVProductionLog uses
[zapcore.EpochTimeEncoder](https://godoc.org/go.uber.org/zap/zapcore#EpochTimeEncoder).
It prints time as floating-point number of seconds since the Unix epoch, and this is a
low cost operation compared to printing time in ISO8609 format.

Other benchmarks (BenchmarkLTSVLog, BenchmarkStandardLog and BenchmarkZapLTSVDevelopmentLog)
uses ISO8609 time format, though there is a slight difference.
BenchmarkZapLTSVDevelopmentLog uses [zapcore.ISO8601TimeEncoder](https://godoc.org/go.uber.org/zap/zapcore#ISO8601TimeEncoder) which prints times with millisecond precision.
The other two prints times with microsecond precision.

```
$ go test -count=10 -bench . -benchmem -cpuprofile=cpu.prof
BenchmarkLTSVLog-2                 	 1000000	      2042 ns/op	      48 B/op	       3 allocs/op
BenchmarkLTSVLog-2                 	 1000000	      2033 ns/op	      48 B/op	       3 allocs/op
BenchmarkLTSVLog-2                 	 1000000	      2073 ns/op	      48 B/op	       3 allocs/op
BenchmarkLTSVLog-2                 	 1000000	      2036 ns/op	      48 B/op	       3 allocs/op
BenchmarkLTSVLog-2                 	 1000000	      2032 ns/op	      48 B/op	       3 allocs/op
BenchmarkLTSVLog-2                 	 1000000	      2018 ns/op	      48 B/op	       3 allocs/op
BenchmarkLTSVLog-2                 	 1000000	      2064 ns/op	      48 B/op	       3 allocs/op
BenchmarkLTSVLog-2                 	 1000000	      2053 ns/op	      48 B/op	       3 allocs/op
BenchmarkLTSVLog-2                 	 1000000	      2045 ns/op	      48 B/op	       3 allocs/op
BenchmarkLTSVLog-2                 	 1000000	      2051 ns/op	      48 B/op	       3 allocs/op
BenchmarkStandardLog-2             	  500000	      2430 ns/op	      96 B/op	       3 allocs/op
BenchmarkStandardLog-2             	  500000	      2384 ns/op	      96 B/op	       3 allocs/op
BenchmarkStandardLog-2             	  500000	      2403 ns/op	      96 B/op	       3 allocs/op
BenchmarkStandardLog-2             	  500000	      2358 ns/op	      96 B/op	       3 allocs/op
BenchmarkStandardLog-2             	  500000	      2470 ns/op	      96 B/op	       3 allocs/op
BenchmarkStandardLog-2             	  500000	      2450 ns/op	      96 B/op	       3 allocs/op
BenchmarkStandardLog-2             	  500000	      2465 ns/op	      96 B/op	       3 allocs/op
BenchmarkStandardLog-2             	  500000	      2370 ns/op	      96 B/op	       3 allocs/op
BenchmarkStandardLog-2             	  500000	      2442 ns/op	      96 B/op	       3 allocs/op
BenchmarkStandardLog-2             	 1000000	      2405 ns/op	      96 B/op	       3 allocs/op
BenchmarkZapLTSVProductionLog-2    	 5000000	       329 ns/op	     128 B/op	       1 allocs/op
BenchmarkZapLTSVProductionLog-2    	 5000000	       348 ns/op	     128 B/op	       1 allocs/op
BenchmarkZapLTSVProductionLog-2    	 5000000	       367 ns/op	     128 B/op	       1 allocs/op
BenchmarkZapLTSVProductionLog-2    	 5000000	       308 ns/op	     128 B/op	       1 allocs/op
BenchmarkZapLTSVProductionLog-2    	 5000000	       364 ns/op	     128 B/op	       1 allocs/op
BenchmarkZapLTSVProductionLog-2    	 3000000	       373 ns/op	     128 B/op	       1 allocs/op
BenchmarkZapLTSVProductionLog-2    	 5000000	       323 ns/op	     128 B/op	       1 allocs/op
BenchmarkZapLTSVProductionLog-2    	 5000000	       368 ns/op	     128 B/op	       1 allocs/op
BenchmarkZapLTSVProductionLog-2    	 5000000	       319 ns/op	     128 B/op	       1 allocs/op
BenchmarkZapLTSVProductionLog-2    	 5000000	       355 ns/op	     128 B/op	       1 allocs/op
BenchmarkZapLTSVDevelopmentLog-2   	  200000	      6089 ns/op	     197 B/op	       4 allocs/op
BenchmarkZapLTSVDevelopmentLog-2   	  200000	      6115 ns/op	     197 B/op	       4 allocs/op
BenchmarkZapLTSVDevelopmentLog-2   	  200000	      6074 ns/op	     197 B/op	       4 allocs/op
BenchmarkZapLTSVDevelopmentLog-2   	  200000	      6172 ns/op	     197 B/op	       4 allocs/op
BenchmarkZapLTSVDevelopmentLog-2   	  200000	      6136 ns/op	     197 B/op	       4 allocs/op
BenchmarkZapLTSVDevelopmentLog-2   	  200000	      6177 ns/op	     197 B/op	       4 allocs/op
BenchmarkZapLTSVDevelopmentLog-2   	  200000	      6107 ns/op	     197 B/op	       4 allocs/op
BenchmarkZapLTSVDevelopmentLog-2   	  200000	      6098 ns/op	     197 B/op	       4 allocs/op
BenchmarkZapLTSVDevelopmentLog-2   	  200000	      6121 ns/op	     197 B/op	       4 allocs/op
BenchmarkZapLTSVDevelopmentLog-2   	  200000	      6229 ns/op	     197 B/op	       4 allocs/op
PASS
ok  	github.com/hnakamur/ltsvlog	67.479s
```

## License
MIT
