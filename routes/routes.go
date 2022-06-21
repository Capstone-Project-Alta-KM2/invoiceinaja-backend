package routes

import (
	"net/http"

	"invoiceinaja/auth"
	"invoiceinaja/domain/user"
	"invoiceinaja/handler"

	"github.com/gin-gonic/gin"
)

func APIRoutes(router *gin.Engine, userHandler *handler.UserHandler, authService auth.Service, userService user.IService) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ping-pong",
		})
	})

	api := router.Group("/api/v1")

	// user
	api.POST("/users", userHandler.UserRegister)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/sessions", userHandler.Login)
	api.POST("/avatars", auth.AuthMiddleware(authService, userService), userHandler.UploadAvatar)
	api.PUT("/users", auth.AuthMiddleware(authService, userService), userHandler.UpdateUser)
}
