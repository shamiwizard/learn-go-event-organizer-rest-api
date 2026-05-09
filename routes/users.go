package routes

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"example.com/event_booking/models"
	"example.com/event_booking/utils"
)

func signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid params", "error": err})
		return
	}
	
	err = user.Save()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create an user", "error": err})
		return
	}

	context.JSON(http.StatusCreated, gin.H{ "message": "User created", "user": user })
}

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid params", "error": err})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid params", "error": err})
		return
	}

	token, err := utils.GenerateJwtToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid params", "error": err})
		return
	}

	fmt.Println(utils.VerifyToken(token))
	
	context.JSON(http.StatusOK, gin.H{ "message": "Sign", "token": token })
}
