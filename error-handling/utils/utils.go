package utils

import (
	"fmt"
	"os"

	"error-handling/errors"
	"error-handling/logger"
)

// HandleError handles an error in a standardized way
// It logs the error with appropriate context and returns an exit code
func HandleError(log *logger.Logger, err error) int {
	if err == nil {
		return 0
	}
	
	// Extract the OpError if available
	opErr, ok := errors.GetOpError(err)
	
	// Log the error with context
	if ok && opErr.Fields != nil {
		log.Error(err.Error(), opErr.Fields)
	} else {
		log.Error(err.Error())
	}
	
	// Log stack trace if available
	stackTrace := errors.GetStackTrace(err)
	if stackTrace != "" {
		log.Debug("Stack trace", stackTrace)
	}
	
	// Determine exit code based on error type
	switch {
	case errors.Is(err, errors.ErrNotFound):
		return 1
	case errors.Is(err, errors.ErrInvalidInput):
		return 2
	case errors.Is(err, errors.ErrUnauthorized):
		return 3
	case errors.Is(err, errors.ErrDatabaseError):
		return 4
	default:
		return 5
	}
}

// PanicOnError panics if the error is not nil
// Useful for errors that should never happen during initialization
func PanicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err))
	}
}

// ExitOnError exits the program with the given status code if the error is not nil
// Useful for fatal errors in command-line applications
func ExitOnError(log *logger.Logger, err error, msg string) {
	if err != nil {
		log.Error(fmt.Sprintf("%s: %v", msg, err))
		os.Exit(1)
	}
}

// MustSucceed is a generic helper that panics on error
// Useful for operations that must not fail
func MustSucceed[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}