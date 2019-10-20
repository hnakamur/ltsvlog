ltsvlog
=======

ltsvlog is a minimalist [LTSV; Labeled Tab-separated Values](http://ltsv.org/) logging library in Go.
See https://godoc.org/github.com/hnakamur/ltsvlog for the API document.

## Warning

* This project is open source, but closed development, no support, no pull request welcome. If you are unsatisfied, feel free to fork it, pick another library, or roll your own.
* I understand this may be a "Don't do that" style in the Go logging and error handling best practices.
* I don't promise the future compatibility. Also this library may not work in the future versions of Go, and a migration path may not be provided, leading to the dead end, users of this library will have to rewrite all of your code. You are on your own with no help.

## An example code and output

An example code:

```
package main

import (
	"github.com/hnakamur/errstack"
	"github.com/hnakamur/ltsvlog"
)

func main() {
	if ltsvlog.Logger.DebugEnabled() {
		ltsvlog.Logger.Debug().String("msg", "This is a debug message").
			String("str", "foo").Int("int", 234).Log()
	}

	ltsvlog.Logger.Info().Fmt("float1", "%3.2f", 3.14).Log()

	if err := outer(); err != nil {
		ltsvlog.Logger.Err(err)
	}
}

func outer() error {
	if err := iInner(); err != nil {
		return errstack.WithLV(errstack.Errorf("add some message here: %s", err), "userID", "user1")
	}
	return nil
}

func iInner() error {
	return errstack.WithLV(errstack.New("some error"), "reqID", "req1")
}
```

An example output:

```
time:2019-10-20T01:18:02.146127Z        level:Debug     msg:This is a debug message     str:foo    int:234
time:2019-10-20T01:18:02.146143Z        level:Info      float1:3.14
time:2019-10-20T01:18:02.146191Z        level:Error     err:add some message here: some error      reqID:req1      userID:user1    stack:main.outer@/home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/main.go:23 main.main@/home/hnakamur/go/src/github.com/hnakamur/ltsvlog/example/main.go:16 runtime.main@/usr/local/go/src/runtime/proc.go:212 runtime.goexit@/usr/local/go/src/runtime/asm_amd64.s:1358
```

Since these log lines are long, please scroll horizontally to the right to see all the output.

## Goals and non-goals

### Goals

* structured logging in LTSV format.
* fast operation and minimum count of memory allocations for Debug and Info.
* log call stack frames with optional labeled values as long as the error message.

### Non-Goals

* compatiblity with the Go logging best practices.
* fast operation and minimum count of memory allocations for Err.
* flexibility and features.

## Benchmark result
[hnakamur/go-log-benchmarks](https://github.com/hnakamur/go-log-benchmarks)

## License
MIT
