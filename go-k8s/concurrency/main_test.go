package main

import (
	"runtime"
	"testing"
)

func BenchmarkAdd(b *testing.B) {
	numbers := generateList(1e7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		add(numbers)
	}
}

func BenchmarkAddConcurrent(b *testing.B) {
	numbers := generateList(1e7)
	goroutines := runtime.GOMAXPROCS(-1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addConcurrent(goroutines, numbers)
	}
}
