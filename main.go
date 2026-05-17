package main

import (
	"example.com/event_booking/db"
	"example.com/event_booking/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	db.InitDB()

	routes.RegisterRoutes(server)

	server.Run(":3000")
}
