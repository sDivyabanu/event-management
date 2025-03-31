package routes

import (
	"backend-challenge/controllers"
	"backend-challenge/middleware"
	"github.com/gin-gonic/gin"
)

func EventRoutes(router *gin.Engine) {
	eventRoutes := router.Group("/api/events")
	eventRoutes.Use(middleware.AuthMiddleware()) 
	{
		eventRoutes.POST("/", controllers.CreateEvent)
		eventRoutes.GET("/", controllers.GetAllEvents)
		eventRoutes.GET("/:id", controllers.GetEventByID)
	}
}
