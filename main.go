package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	server.GET("/events", getEvents)

	server.Run()
}

func getEvents(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "hellow", "sec_message": 2})
}
