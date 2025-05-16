package main

import (
	"fmt"
	"net/http"

	"os"
	"runtime/pprof"

	_ "net/http/pprof"
)

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func main() {
	// Start pprof HTTP server for live profiling
	go func() {
		fmt.Println("pprof server listening on :6060")
		http.ListenAndServe(":6060", nil)
	}()

	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("could not create CPU profile:", err)
		return
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Run a CPU-intensive task
	fmt.Println("Calculating fib(35)...")
	result := fib(35)
	fmt.Println("fib(35) =", result)

	// Prevent main from exiting so the pprof server stays up
	fmt.Println("Press Ctrl+C to exit and stop the pprof server.")
	// select {} // Block forever
}
