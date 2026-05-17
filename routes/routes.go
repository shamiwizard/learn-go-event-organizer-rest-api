package routes

import (
	"example.com/event_booking/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	v1 := server.Group("/api/v1/")
	{

		v1.GET("/events", getEvents)
		v1.GET("/events/:id", getEvent)

		authGroup := v1.Group("/")
		authGroup.Use(middleware.Authenticate)

		authGroup.POST("/events", createEvent)
		authGroup.PUT("/events/:id", updateEvent)
		authGroup.DELETE("/events/:id", deleteEvent)
		authGroup.POST("/events/:id/register", registerForEvent)
		authGroup.DELETE("/events/:id/register", cancelRegistration)

		v1.POST("/signup", signup)
		v1.POST("/login", login)
	}
}
