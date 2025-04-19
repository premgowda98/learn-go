package main

import (
	"slogtest/internal/pkg/logger"
	"slogtest/internal/pkg/service"
)

func main() {

	kc_log := logger.NewLogger("main")

	// Create a new user service
	userService := service.NewUserService()

	// Create some sample users
	user := service.User{
		ID:       "1",
		Username: "john_doe",
		Email:    "john@example.com",
	}
	// Process users and demonstrate different log levels
	kc_log.Info("starting user processing")

	err := userService.CreateUser(user)

	if err != nil {
		kc_log.Error("failed to create user",
			"error", err,
			"user_id", user.ID,
		)
	}

	kc_log.Info("application finished")
}
