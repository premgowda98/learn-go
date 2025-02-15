package handlers

import (
	"database/sql"
	"net/http"
	"project/user-management/internal/models"
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

func UpdateUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err !=nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		user, err := services.GetUserByID(db, id)

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		if err !=nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong"})
			return
		}

		var userUpdate models.User

		if err := c.BindJSON(&userUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err = services.UpdateUser(db, &userUpdate, id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
