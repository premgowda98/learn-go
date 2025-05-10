package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Standard error types that can be checked with errors.Is
var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInternal      = errors.New("internal error")
	ErrDatabaseError = errors.New("database error")
)

// StackTracer is an interface for errors that can provide a stack trace
type StackTracer interface {
	StackTrace() string
}

// OpError represents an operational error with context
type OpError struct {
	// Op is the operation that failed (e.g., "db.query", "http.request")
	Op string
	// Err is the underlying error
	Err error
	// Message is a human-readable error message
	Message string
	// Fields contains additional structured context for the error
	Fields map[string]interface{}
	// stack contains the stack trace
	stack string
}

// Error implements the error interface
func (e *OpError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

// Unwrap allows errors.Is and errors.As to work with our custom error
func (e *OpError) Unwrap() error {
	return e.Err
}

// StackTrace returns the stack trace for this error
func (e *OpError) StackTrace() string {
	return e.stack
}

// captureStack captures the current stack trace, skipping the given number of frames
func captureStack(skip int) string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	
	var builder strings.Builder
	for {
		frame, more := frames.Next()
		
		// Skip runtime frames
		if !strings.Contains(frame.File, "runtime/") {
			fmt.Fprintf(&builder, "%s:%d - %s\n", frame.File, frame.Line, frame.Function)
		}
		
		if !more {
			break
		}
	}
	
	return builder.String()
}

// E creates a new OpError with the given operation and underlying error
func E(op string, err error, opts ...func(*OpError)) error {
	if err == nil {
		return nil
	}
	
	// Create a new OpError
	e := &OpError{
		Op:    op,
		Err:   err,
		stack: captureStack(3), // Skip 3 stack frames to get to the caller's caller
	}
	
	// Apply the options
	for _, opt := range opts {
		opt(e)
	}
	
	return e
}

// WithMessage adds a message to the error
func WithMessage(message string) func(*OpError) {
	return func(e *OpError) {
		e.Message = message
	}
}

// WithField adds a field to the error
func WithField(key string, value interface{}) func(*OpError) {
	return func(e *OpError) {
		if e.Fields == nil {
			e.Fields = make(map[string]interface{})
		}
		e.Fields[key] = value
	}
}

// WithFields adds multiple fields to the error
func WithFields(fields map[string]interface{}) func(*OpError) {
	return func(e *OpError) {
		if e.Fields == nil {
			e.Fields = make(map[string]interface{})
		}
		for k, v := range fields {
			e.Fields[k] = v
		}
	}
}

// GetOpError extracts an OpError from the error chain
func GetOpError(err error) (*OpError, bool) {
	var opErr *OpError
	if errors.As(err, &opErr) {
		return opErr, true
	}
	return nil, false
}

// GetStackTrace extracts stack trace from the error chain if available
func GetStackTrace(err error) string {
	var stackTracer StackTracer
	if errors.As(err, &stackTracer) {
		return stackTracer.StackTrace()
	}
	return ""
}

// Is checks if the target error is in the error chain
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}