package routes

import (
	"github.com/gin-gonic/gin"
	"example.com/event_booking/middleware"
)


func RegisterRoutes(server *gin.Engine) {
	v1 := server.Group("/api/v1/")
	{
		v1.GET("/events", getEvents)
		v1.GET("/events/:id", getEvent)
		v1.POST("/events", middleware.Authorize(), createEvent)
		v1.PUT("/events/:id", middleware.Authorize(), updateEvent)
		v1.DELETE("/events/:id",middleware.Authorize(), deleteEvent)
		v1.POST("/signup",signup)
		v1.POST("/login", login)
	}
}

