package routes

import (
	"database/sql"
	"project/user-management/internal/handlers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/login", handlers.LoginHandler(db))
	r.POST("/register", handlers.RegisterHandler(db))
	r.GET("/users/:id", handlers.GetUserHander(db))
}
