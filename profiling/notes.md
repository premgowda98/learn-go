# Go Profiling: Complete Guide

[Video](https://youtu.be/R_C2KqmlqY8)

Profiling is the process of measuring the performance of your Go programs to identify bottlenecks, excessive memory usage, or inefficient code paths. Go provides built-in tools for CPU, memory, and other types of profiling.

## Why Profile?
- Find slow or inefficient code
- Reduce CPU and memory usage
- Improve application performance

## Types of Profiling in Go
- **CPU Profiling**: Shows where your program spends its time (which functions use the most CPU).
- **Memory (Heap) Profiling**: Shows memory allocation patterns and which parts of your code use the most memory.
- **Block Profiling**: Shows where goroutines are waiting (e.g., on locks or channel operations).
- **Mutex Profiling**: Shows contention on mutexes.
- **Goroutine Profiling**: Shows stack traces of all goroutines.

---

## CPU Profiling in Go: Two Main Approaches

### 1. Manual Profiling with `runtime/pprof.StartCPUProfile`
This approach lets you programmatically start and stop CPU profiling and save the results to a file.

**Example:**
```go
import (
    "os"
    "runtime/pprof"
)

func main() {
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    // ...your code to profile...
}
```
- Run your program to generate `cpu.prof`.
- Analyze with:
  ```sh
  go tool pprof main cpu.prof
  ```
- Use commands like `top`, `list`, `web` inside pprof for analysis.

### 2. Live Profiling with `net/http/pprof`
This approach exposes a web server with profiling endpoints for live inspection and profile download.

**Example:**
```go
import _ "net/http/pprof"
import "net/http"

func main() {
    go func() {
        http.ListenAndServe(":6060", nil)
    }()
    // ...your code...
}
```
- Run your program, then visit [http://localhost:6060/debug/pprof/](http://localhost:6060/debug/pprof/) in your browser.
- Download profiles (e.g., `/debug/pprof/profile` for CPU).
- You can fetch and analyze a profile directly with:
  ```sh
  go tool pprof http://localhost:6060/debug/pprof/profile
  ```
  - By default, this collects 30 seconds of CPU profile data. You can change the duration:
    ```sh
    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10
    ```
- After fetching, use pprof commands (`top`, `list`, `web`) for analysis.
- For a web UI:
  ```sh
  go tool pprof -http=:8080 main cpu.prof
  ```

---

## When to Use Each Approach
- **Manual profiling** is useful for short-lived programs or when you want to profile a specific code section.
- **net/http/pprof** is ideal for long-running services, allowing live inspection and remote profile collection.

## Notes on Profile Collection
- The default duration for CPU profile collection via `/debug/pprof/profile` is 30 seconds. You can set a custom duration with the `seconds` query parameter.
- The profile is loaded into pprof for analysis; save it manually if you want a `.prof` file.

## Common pprof Endpoints
- `/debug/pprof/profile` — CPU profile (default 30s)
- `/debug/pprof/heap` — Heap profile
- `/debug/pprof/block` — Block profile
- `/debug/pprof/mutex` — Mutex profile
- `/debug/pprof/goroutine` — Goroutine dump

## Visualizing Profiles
- Use `go tool pprof -http=:8080 main cpu.prof` to open a web UI for interactive exploration.
- Use `web` command in pprof to generate SVG call graphs (requires Graphviz).

## Best Practices
- Profile in production-like environments for realistic data.
- Profile both CPU and memory for a complete picture.
- Use flamegraphs for visualizing hotspots.
- Remove profiling code from production builds if not needed.

## Useful Tools
- [pprof](https://pkg.go.dev/net/http/pprof) (standard)
- [go-torch](https://github.com/uber-archive/go-torch) (flamegraphs)
- [speedscope](https://www.speedscope.app/) (visualization)

---
Profiling is essential for writing high-performance Go applications. Use the built-in tools to find and fix bottlenecks efficiently.