package handlers

import (
	"database/sql"
	"net/http"
	"project/user-management/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserHander(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err !=nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		user, err := services.GetUserByID(db, id)

		if err !=nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong"})
			return
		}

		c.JSON(http.StatusOK, user)

	}
}
