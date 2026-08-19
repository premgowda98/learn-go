package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	// Example 1: Context with timeout
	fmt.Println("=== Context with Timeout Example ===")
	timeoutExample()

	// Example 2: Context with cancel
	fmt.Println("\n=== Context with Cancel Example ===")
	cancelExample()

	// Example 3: Context with values
	fmt.Println("\n=== Context with Values Example ===")
	valueExample()

	// Example 4: Context with deadline
	fmt.Println("\n=== Context with Deadline Example ===")
	deadlineExample()
}

func timeoutExample() {
	// Create a context that will timeout after 2 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Always remember to cancel to release resources

	// Simulate some work
	go func() {
		time.Sleep(3 * time.Second) // This work takes 3 seconds
		fmt.Println("Work done (but too late)")
	}()

	// Wait for the work to complete or timeout
	select {
	case <-ctx.Done():
		fmt.Println("Operation timed out:", ctx.Err())
	}
}

func cancelExample() {
	ctx, cancel := context.WithCancel(context.Background())

	// Start a goroutine that watches for cancellation
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine cancelled:", ctx.Err())
				return
			default:
				fmt.Println("Doing work...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// Let it work for 1.5 seconds
	time.Sleep(1500 * time.Millisecond)
	cancel()                           // Cancel the operation
	time.Sleep(100 * time.Millisecond) // Wait to see the cancellation effect
}

func valueExample() {
	// Create a context with some values
	ctx := context.Background()
	ctx = context.WithValue(ctx, "userID", "123")
	ctx = context.WithValue(ctx, "authToken", "xyz789")

	// Simulate passing context through function calls
	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	// Retrieve values from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		log.Println("userID not found in context")
		return
	}

	authToken, ok := ctx.Value("authToken").(string)
	if !ok {
		log.Println("authToken not found in context")
		return
	}

	fmt.Printf("Processing request for user %s with token %s\n", userID, authToken)
}

func deadlineExample() {
	// Create a context with a deadline of 1 second from now
	deadline := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// Try to do some work before the deadline
	select {
	case <-time.After(2 * time.Second):
		fmt.Println("Work completed (this won't be printed)")
	case <-ctx.Done():
		fmt.Println("Work cancelled due to deadline:", ctx.Err())
	}
}
