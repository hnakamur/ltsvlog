ltsvlog
=======

ltsvlog is a minimalist [LTSV; Labeled Tab-separated Values](http://ltsv.org/) logging library in Go.
See https://godoc.org/github.com/hnakamur/ltsvlog for the API document.

## An example code and output

See [an example code](https://github.com/hnakamur/ltsvlog/blob/master/cmd/example/main.go)

```
$ go run cmd/example/main.go
time:2016-05-23T15:57:17.748330442Z     level:Debug     msg:This is a debug message     key:key1        intValue:234
time:2016-05-23T15:57:17.748360741Z     level:Info      msg:hello, world        key:key1        value:value1
time:2016-05-23T15:57:17.748366628Z     level:Info      msg:goodbye, world      foo:bar nilValue:<nil>  bytes:0x612f62
time:2016-05-23T15:57:17.748420201Z     level:Debug     msg:stack trace example stack:goroutine 1 [running]: [main.b() /root/gocode/src/github.com/hnakamur/ltsvlog/cmd/example/main.go:26 +0x49] [main.a() /root/gocode/src/github.com/hnakamur/ltsvlog/cmd/example/main.go:22 +0x14] [main.main() /root/gocode/src/github.com/hnakamur/ltsvlog/cmd/example/main.go:18 +0xa59]
```

## Benchmark result

```
$ go test -bench . -benchmem
testing: warning: no tests to run
PASS
BenchmarkLTSVLog-2       1000000              1242 ns/op             245 B/op          3 allocs/op
BenchmarkStandardLog-2   1000000              1186 ns/op             235 B/op          3 allocs/op
ok      github.com/hnakamur/ltsvlog     2.486s
```

## License
MIT
