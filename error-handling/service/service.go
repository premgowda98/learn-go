package service

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"error-handling/errors"
)

// UserService represents a mock user service
type UserService struct{}

// User represents a user entity
type User struct {
	ID       int
	Username string
	Email    string
}

// GetUserByID retrieves a user by ID and demonstrates different error handling scenarios
func (s *UserService) GetUserByID(id int) (*User, error) {
	const op = "service.GetUserByID"
	
	// Validate input
	if id <= 0 {
		return nil, errors.E(op, errors.ErrInvalidInput,
			errors.WithMessage("user ID must be positive"),
			errors.WithField("id", id))
	}
	
	// Mock DB operation that might fail
	user, err := s.findUserInDB(id)
	if err != nil {
		return nil, errors.E(op, err,
			errors.WithMessage(fmt.Sprintf("failed to get user with ID %d", id)),
			errors.WithField("id", id))
	}
	
	return user, nil
}

// CreateUser demonstrates error handling in a create operation
func (s *UserService) CreateUser(username, email string) (*User, error) {
	const op = "service.CreateUser"
	
	// Validate input
	if username == "" {
		return nil, errors.E(op, errors.ErrInvalidInput,
			errors.WithMessage("username cannot be empty"))
	}
	
	if email == "" {
		return nil, errors.E(op, errors.ErrInvalidInput,
			errors.WithMessage("email cannot be empty"))
	}
	
	// Check if user already exists
	exists, err := s.checkUserExists(username)
	if err != nil {
		return nil, errors.E(op, err,
			errors.WithMessage("failed to check if user exists"),
			errors.WithField("username", username))
	}
	
	if exists {
		return nil, errors.E(op, errors.ErrInvalidInput,
			errors.WithMessage("username already exists"),
			errors.WithField("username", username))
	}
	
	// Mock saving user to database
	user, err := s.saveUserToDB(username, email)
	if err != nil {
		return nil, errors.E(op, err,
			errors.WithMessage("failed to save user to database"),
			errors.WithFields(map[string]interface{}{
				"username": username,
				"email":    email,
			}))
	}
	
	return user, nil
}

// UpdateUserEmail demonstrates error handling in an update operation
func (s *UserService) UpdateUserEmail(id int, newEmail string) error {
	const op = "service.UpdateUserEmail"
	
	// Validate input
	if id <= 0 {
		return errors.E(op, errors.ErrInvalidInput,
			errors.WithMessage("user ID must be positive"),
			errors.WithField("id", id))
	}
	
	if newEmail == "" {
		return errors.E(op, errors.ErrInvalidInput,
			errors.WithMessage("email cannot be empty"))
	}
	
	// Check if user exists
	_, err := s.findUserInDB(id)
	if err != nil {
		return errors.E(op, err,
			errors.WithMessage(fmt.Sprintf("failed to find user with ID %d", id)),
			errors.WithField("id", id))
	}
	
	// Mock update operation
	err = s.updateUserEmailInDB(id, newEmail)
	if err != nil {
		return errors.E(op, err,
			errors.WithMessage("failed to update user email"),
			errors.WithFields(map[string]interface{}{
				"id":       id,
				"newEmail": newEmail,
			}))
	}
	
	return nil
}

// Mock helper functions that simulate database operations

func (s *UserService) findUserInDB(id int) (*User, error) {
	// For demonstration, generate random errors
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(10)
	
	switch {
	case r < 3:
		// Simulate user not found
		return nil, errors.E("db.findUser", errors.ErrNotFound,
			errors.WithMessage(fmt.Sprintf("user with ID %d not found", id)))
	case r < 5:
		// Simulate database error
		return nil, errors.E("db.findUser", errors.ErrDatabaseError,
			errors.WithMessage("connection timeout"),
			errors.WithField("attempt", 3))
	case r < 7:
		// Simulate SQL error
		return nil, errors.E("db.findUser", sql.ErrNoRows)
	}
	
	// Success case
	return &User{
		ID:       id,
		Username: fmt.Sprintf("user%d", id),
		Email:    fmt.Sprintf("user%d@example.com", id),
	}, nil
}

func (s *UserService) checkUserExists(username string) (bool, error) {
	// For demonstration, generate random errors
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(10)
	
	if r < 2 {
		// Simulate database error
		return false, errors.E("db.checkUserExists", errors.ErrDatabaseError,
			errors.WithMessage("failed to execute query"))
	}
	
	// Simulate existing user
	return username == "admin" || r < 4, nil
}

func (s *UserService) saveUserToDB(username, email string) (*User, error) {
	// For demonstration, generate random errors
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(10)
	
	if r < 3 {
		// Simulate database error
		return nil, errors.E("db.saveUser", errors.ErrDatabaseError,
			errors.WithMessage("failed to insert new user"))
	}
	
	// Success case
	return &User{
		ID:       rand.Intn(1000) + 1,
		Username: username,
		Email:    email,
	}, nil
}

func (s *UserService) updateUserEmailInDB(id int, email string) error {
	// For demonstration, generate random errors
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(10)
	
	if r < 3 {
		// Simulate database error
		return errors.E("db.updateUserEmail", errors.ErrDatabaseError,
			errors.WithMessage("failed to update user"))
	}
	
	return nil
}