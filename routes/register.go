package routes

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"example.com/event_booking/models"
)

func registerForEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not procces a request." })
		return
	}

	event, err := models.FindEvent(id)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{})
		return
	}

	user := context.MustGet("currentUser").(models.User)
	err = event.RegisterUser(&user)


	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not process a request.", "error": fmt.Sprint(err) })
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User registered to an event" })
}

func cancelRegistration(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not procces a request." })
		return
	}

	event, err := models.FindEvent(id)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{})
		return
	}

	user := context.MustGet("currentUser").(models.User)
	err = event.CancelRegistration(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not process a request.", "error": fmt.Sprint(err) })
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Registration canceled" })
}

