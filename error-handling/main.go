package main

import (
	"fmt"
	"os"
	"strconv"

	"error-handling/errors"
	"error-handling/logger"
	"error-handling/service"
	"error-handling/utils"
)

func main() {
	// Initialize logger
	log := logger.New()
	
	// Create an instance of UserService
	userService := &service.UserService{}
	
	// Demonstrate different error scenarios
	fmt.Println("=== Error Handling Demonstration ===")
	fmt.Println()
	
	// Case 1: Invalid input - negative ID
	fmt.Println("Case 1: Invalid input - negative ID")
	_, err := userService.GetUserByID(-1)
	if err != nil {
		printErrorDetails(log, err)
	}
	fmt.Println()
	
	// Case 2: Get user by ID (might succeed or fail with different errors)
	fmt.Println("Case 2: Get user by ID (random result)")
	user, err := userService.GetUserByID(42)
	if err != nil {
		printErrorDetails(log, err)
	} else {
		log.Info("User found", user)
	}
	fmt.Println()
	
	// Case 3: Creating a user with empty fields
	fmt.Println("Case 3: Creating a user with empty fields")
	_, err = userService.CreateUser("", "")
	if err != nil {
		printErrorDetails(log, err)
	}
	fmt.Println()
	
	// Case 4: Create a user (might succeed or fail)
	fmt.Println("Case 4: Create a user (random result)")
	newUser, err := userService.CreateUser("testuser", "test@example.com")
	if err != nil {
		printErrorDetails(log, err)
	} else {
		log.Info("User created", newUser)
	}
	fmt.Println()
	
	// Case 5: Update a user's email (might succeed or fail)
	fmt.Println("Case 5: Update a user's email (random result)")
	err = userService.UpdateUserEmail(123, "new@example.com")
	if err != nil {
		printErrorDetails(log, err)
	} else {
		log.Info("User email updated successfully")
	}
	fmt.Println()
	
	// Case 6: Demonstrate command-line arguments handling
	fmt.Println("Case 6: Process command-line arguments (if provided)")
	if len(os.Args) > 1 {
		// Try to process first argument as user ID
		processCommandLineArgs(log, userService, os.Args[1])
	}
}

// Helper function to demonstrate error handling with command line args
func processCommandLineArgs(log *logger.Logger, userService *service.UserService, idArg string) {
	const op = "main.processCommandLineArgs"
	
	// Try to convert the argument to an integer
	id, err := strconv.Atoi(idArg)
	if err != nil {
		err = errors.E(op, errors.ErrInvalidInput,
			errors.WithMessage("failed to parse user ID"),
			errors.WithField("input", idArg))
		printErrorDetails(log, err)
		return
	}
	
	// Try to get the user by ID
	user, err := userService.GetUserByID(id)
	if err != nil {
		printErrorDetails(log, err)
		return
	}
	
	log.Info("User found from command line argument", user)
}

// Helper function to print error details
func printErrorDetails(log *logger.Logger, err error) {
	fmt.Println("❌ Error occurred:")
	fmt.Printf("   Message: %v\n", err)
	
	// Check error type
	switch {
	case errors.Is(err, errors.ErrNotFound):
		fmt.Println("   Type: Resource Not Found Error")
	case errors.Is(err, errors.ErrInvalidInput):
		fmt.Println("   Type: Invalid Input Error")
	case errors.Is(err, errors.ErrDatabaseError):
		fmt.Println("   Type: Database Error")
	default:
		fmt.Println("   Type: Unknown Error")
	}
	
	// Extract and print stack trace
	stackTrace := errors.GetStackTrace(err)
	if stackTrace != "" {
		fmt.Println("   Stack Trace Available in Logs")
	}
	
	// Use our error handler (would handle logging and returning exit code)
	exitCode := utils.HandleError(log, err)
	fmt.Printf("   Exit Code: %d\n", exitCode)
}