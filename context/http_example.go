package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	Data string `json:"data"`
}

func slowOperation(ctx context.Context) (string, error) {
	// Create a channel for our result
	resultCh := make(chan string)

	go func() {
		// Simulate a slow operation (e.g., database query or external API call)
		time.Sleep(2 * time.Second)
		resultCh <- "Operation completed successfully"
	}()

	// Wait for either the operation to complete or context to be cancelled
	select {
	case result := <-resultCh:
		return result, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Create a context with timeout for this request
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Add some values to the context (simulating middleware adding auth info)
	ctx = context.WithValue(ctx, "requestID", time.Now().UnixNano())

	// Perform the slow operation
	result, err := slowOperation(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	// Get the request ID from context
	requestID := ctx.Value("requestID").(int64)

	// Prepare and send response
	response := Response{
		Data: fmt.Sprintf("RequestID: %d, Result: %s", requestID, result),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func startServer() {
	http.HandleFunc("/slowop", handleRequest)
	fmt.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", nil)
}