package main

import (
	"github.com/gin-gonic/gin"
	"example.com/event_booking/routes"
	"example.com/event_booking/db"
)

func main() {
	server := gin.Default()
	db.InitDB()

	routes.RegisterRoutes(server)

	server.Run(":3000")
}

