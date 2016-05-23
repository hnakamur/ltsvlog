ltsvlog
=======

ltsvlog is a minimalist [LTSV; Labeled Tab-separated Values](http://ltsv.org/) logging library in Go.
See https://godoc.org/github.com/hnakamur/ltsvlog for the API document.

## Benchmark result

```
$ go test -bench . -benchmem
testing: warning: no tests to run
PASS
BenchmarkLTSVLog-2       1000000              1225 ns/op             245 B/op          3 allocs/op
BenchmarkStandardLog-2   1000000              1223 ns/op             235 B/op          3 allocs/op
ok      github.com/hnakamur/ltsvlog     2.512s
```

## License
MIT
