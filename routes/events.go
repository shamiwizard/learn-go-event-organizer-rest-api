package routes

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"example.com/event_booking/models"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not find events. Try later." })

		return
	}

	context.JSON(http.StatusOK, gin.H{ "events": events })
}

func getEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not procces an ID." })
		return
	}

	event, err := models.FindEvent(id)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{})
		return
	}

	
	context.JSON(http.StatusOK, gin.H{"event": event})
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create an event", "error": err})
		return
	}

	event.UserID = context.MustGet("currentUser").(models.User).ID
	err = event.Save()

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Could not process the data." })
	} else {
		context.JSON(http.StatusCreated, gin.H{ "message": "Event created", "event": event})
	}
}

func updateEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not procces an ID." })
		return
	}
	
	event, err := models.FindEvent(id)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{})
		return
	}

	err = context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not update an event", "error": err})
		return
	}

	err = event.Update()


	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not update an event", "error": err})
		return
	}

	context.JSON(http.StatusOK, gin.H{ "message": "Event updated", "event": event})
}

func deleteEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not procces an ID." })
		return
	}


	event, err := models.FindEvent(id)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{})
		return
	}

	err = event.Delete()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not delete an event", "error": err})
		return
	}

	context.JSON(http.StatusOK, gin.H{ "message": "Event deleted" })
}
