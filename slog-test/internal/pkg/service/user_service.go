package service

import (
	"fmt"
	"log/slog"
	"slogtest/internal/pkg/logger"
)

var kc_log *slog.Logger = logger.GetLogger("user_service")

type User struct {
	ID       string
	Username string
	Email    string
}

type UserService struct {
	// Remove logger field as we'll use the global logger directly
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(user User) error {

	kc_log.Info("creating new user")

	if user.Username == "" {
		slog.Error("failed to create user: username is required")
		return fmt.Errorf("username is required")
	}

	if user.Email == "" {
		slog.Error("failed to create user: email is required")
		return fmt.Errorf("email is required")
	}

	// Simulate some work
	kc_log.Debug("validating user data")

	kc_log.Error("failed here")

	// Simulate successful user creation
	kc_log.Info("user created successfully")
	return fmt.Errorf("something went wrong")
}
