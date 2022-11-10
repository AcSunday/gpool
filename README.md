## gpool is an ants based pool that automatically resizes goroutines.

### 

```
➜  gpool git:(main) go test -bench . -benchmem -benchtime=5s -timeout=30s
goos: darwin
goarch: amd64
pkg: gpool
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkGoPool-12                 10000            781941 ns/op               1 B/op          0 allocs/op
BenchmarkGoroutine-12              10000            961333 ns/op             422 B/op          1 allocs/op
PASS
ok      gpool   17.562s

```