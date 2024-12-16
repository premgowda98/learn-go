package routes

import (
	"net/http"
	"project/restapi/models"
	"project/restapi/utils"

	"github.com/gin-gonic/gin"
)

func singup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Some fields are missing"})
		return
	}

	err = user.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save user"})
		return 
	}

	context.JSON(http.StatusCreated, gin.H{"message": "EveUsernt Created"})
}


func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Some fields are missing"})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to validate user"})
		return 
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to validate user"})
		return 
	}

	context.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token})
}