package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	v1 := server.Group("/api/v1/")
	{
		v1.GET("/events", getEvents)
		v1.GET("/events/:id", getEvent)
		v1.POST("/events", createEvent)
	}
}
