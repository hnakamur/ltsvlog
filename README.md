ltsvlog
=======

ltsvlog is a minimalist [LTSV; Labeled Tab-separated Values](http://ltsv.org/) logging library in Go.
See https://godoc.org/github.com/hnakamur/ltsvlog for the API document.

## Benchmark result

```
# go test -bench . -benchmem -cpuprofile -memprofile
...(snip)...
BenchmarkLTSVLog-2       2000000               962 ns/op             197 B/op          0 allocs/op
BenchmarkStandardLog-2   1000000              1194 ns/op             235 B/op          3 allocs/op
...(snip)...
```

## License
MIT
