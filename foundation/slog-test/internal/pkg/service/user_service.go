package service

import (
	"fmt"
	"log/slog"
)

type User struct {
	ID       string
	Username string
	Email    string
}

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(user User) error {

	slog.Info("creating new user")

	if user.Username == "" {
		slog.Error("failed to create user: username is required")
		return fmt.Errorf("username is required")
	}

	if user.Email == "" {
		slog.Error("failed to create user: email is required")
		return fmt.Errorf("email is required")
	}

	slog.Debug("validating user data")

	slog.Error("failed here")

	slog.Info("user created successfully", "user", 123)
	return fmt.Errorf("something went wrong")
}
