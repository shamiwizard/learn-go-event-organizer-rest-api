package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"example.com/event_booking/models"
	"example.com/event_booking/db"
)

func main() {
	server := gin.Default()
	db.InitDB()

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})


	server.GET("/events", getEvents)
	server.POST("/events", createEvent)

	server.Run(":3000")
}

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not find events. Try later." })

		return
	}

	context.JSON(http.StatusOK, gin.H{ "events": events })
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create an event", "error": err})
		return
	}

	event.UserID = 1
	err = event.Save()

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Could not process the data." })
	} else {
		context.JSON(http.StatusCreated, gin.H{ "message": "Event created", "event": event})
	}
}
