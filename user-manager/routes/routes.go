package routes

import (
	"github.com/gin-gonic/gin"
	"user-manager/handlers"
)

func RegisterRoutes(router *gin.Engine) {
	user := router.Group("/users")
	{
		user.POST("/", handlers.CreateUser)
		user.GET("/", handlers.GetUsers)
		user.GET("/:id", handlers.GetUserByID)
		user.PUT("/:id", handlers.UpdateUser)
		user.DELETE("/:id", handlers.DeleteUser)
	}
}
