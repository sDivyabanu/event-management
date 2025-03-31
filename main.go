package main

import (
	"backend-challenge/config"
	"backend-challenge/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	router := gin.Default()

	routes.UserRoutes(router)
	routes.EventRoutes(router)

	log.Println("🚀 Server running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}