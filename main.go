package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"example.com/event_booking/models"

)

func main() {
	server := gin.Default()

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
	events := models.GetAllEvents()
	context.JSON(http.StatusOK, gin.H{"message": "hellow", "events": events })
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create an event", "error": err})
		return
	}

	event.ID = 1
	event.UserID = 1
	event.Save()

	context.JSON(http.StatusCreated, gin.H{ "message": "Event created", "event": event})
}
