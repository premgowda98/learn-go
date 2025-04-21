package main

import (
	"log/slog"
	"slogtest/internal/pkg/logger"
	"slogtest/internal/pkg/service"
)

func main() {
	logger.NewLogger("main")

	userService := service.NewUserService()

	user := service.User{
		ID:       "1",
		Username: "john_doe",
		Email:    "john@example.com",
	}
	slog.Info("starting user processing")

	err := userService.CreateUser(user)

	if err != nil {
		slog.Error("failed to create user",
			"error", err,
			"user_id", user.ID,
		)
	}

	slog.Info("application finished")
}
