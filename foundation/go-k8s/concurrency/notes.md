# Concurrency and Parallelism in Go

### With Single OS Thread

```bash
GOGC=off go test -cpu 1 -run none -bench . -benchtime 3s

# results
goos: linux
goarch: amd64
pkg: concurrency-parallelism
cpu: Intel(R) Core(TM) i7-7600U CPU @ 2.80GHz
BenchmarkAdd                1358           2655689 ns/op
BenchmarkAddConcurrent       573           6308321 ns/op
PASS
ok      concurrency-parallelism 9.129s
```

### With 4 CPU & OS Threads

```bash
GOGC=off go test -cpu 4 -run none -bench . -benchtime 3s

# results
goos: linux
goarch: amd64
pkg: concurrency-parallelism
cpu: Intel(R) Core(TM) i7-7600U CPU @ 2.80GHz
BenchmarkAdd-4                      1348           2657827 ns/op
BenchmarkAddConcurrent-4             697           5147838 ns/op
PASS
ok      concurrency-parallelism 8.947s
```


```bash
GOGC=off go test -cpu 1,4 -run none -bench . -benchtime 3s

goos: linux
goarch: amd64
pkg: concurrency-parallelism
cpu: Intel(R) Core(TM) i7-7600U CPU @ 2.80GHz
BenchmarkAdd                        1347           2728243 ns/op
BenchmarkAdd-4                      1323           2651952 ns/op
BenchmarkAddConcurrent               562           6357607 ns/op
BenchmarkAddConcurrent-4             787           4621085 ns/op
PASS
ok      concurrency-parallelism 17.924s
```