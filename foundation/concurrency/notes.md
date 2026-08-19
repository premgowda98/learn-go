## Concurrency

1. Go has a different approach for concurrency
    1. Don't communicate by sharing memory, share memory by communicating

## GoRoutines

1. Runs in background

## Race Condition

1. Mutex -> to deal with race condition
2. Occurs when 2 goroutine access same resource
3. To check for race `go run -race .`
