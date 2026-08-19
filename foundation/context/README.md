# Go Context Package Examples

This directory contains comprehensive examples demonstrating the usage of Go's `context` package. The context package is primarily used for:
- Cancellation signals
- Deadlines
- Passing request-scoped values
- Timeouts

## Key Concepts

1. **Context Types**:
   - `context.Background()`: Root context, typically used in main function or program initialization
   - `context.TODO()`: Placeholder when you're unsure which context to use
   - `context.WithCancel()`: Creates a cancellable context
   - `context.WithDeadline()`: Creates a context that will be cancelled at a specific time
   - `context.WithTimeout()`: Creates a context that will be cancelled after a duration
   - `context.WithValue()`: Creates a context with a key-value pair

2. **Best Practices**:
   - Always pass context as the first parameter to functions
   - Don't store contexts inside structs
   - Use context values only for request-scoped data
   - Always call cancel when you're done with a context
   - Don't pass nil contexts, use context.TODO() if you're unsure

## Examples in this Directory

1. **main.go**: Contains basic examples of:
   - Context with timeout
   - Context with cancel
   - Context with values
   - Context with deadline

2. **http_example.go**: Demonstrates practical usage in HTTP servers:
   - Request timeout handling
   - Request-scoped values
   - Cancellation propagation
   - Integration with http.Request.Context()

## Running the Examples

1. Basic examples:
```bash
go run main.go
```

2. HTTP server example:
```bash
# Run the server
go run http_example.go

# In another terminal, test with curl:
curl http://localhost:8080/slowop
```

## Common Use Cases

1. **HTTP Servers**: Managing request timeouts and cancellations
2. **Database Operations**: Setting query timeouts
3. **API Calls**: Managing timeouts for external service calls
4. **Resource Cleanup**: Ensuring proper cleanup when operations are cancelled
5. **Request Tracing**: Passing request-specific data through call stack