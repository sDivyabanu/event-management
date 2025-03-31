package routes

import (
	"backend-challenge/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/api/users")
	{
		userRoutes.POST("/signup", controllers.SignUp)
		userRoutes.POST("/login", controllers.Login)
	}
}
