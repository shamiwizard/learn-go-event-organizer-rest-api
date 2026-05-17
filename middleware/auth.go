package middleware

import (
	"example.com/event_booking/models"
	"example.com/event_booking/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authenticate(context *gin.Context) {
	authToken := context.GetHeader("Authorization")

	if len(authToken) > 7 && authToken[:7] == "Bearer " {
		authToken = authToken[7:]
	}

	claims, err := utils.VerifyToken(authToken)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	id, idOk := claims["userId"].(float64)
	email, emailOk := claims["userEmail"].(string)

	if !(idOk || emailOk) {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	user := models.FindByIdEmail(int64(id), email)

	if user == nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	context.Set("currentUser", *user)

	context.Next()
}
