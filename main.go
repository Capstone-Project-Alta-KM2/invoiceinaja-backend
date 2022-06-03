package main

import (
	"invoiceinaja/auth"
	"invoiceinaja/config"
	"invoiceinaja/database"
	"invoiceinaja/domain/user"
	"invoiceinaja/handler"

	"github.com/gin-gonic/gin"
	"github.com/treblle/treblle-go"
)

func main() {
	treblle.Configure(treblle.Configuration{
		ProjectID: "b6lOPdqzPnWVRR6Z",
	})

	conf := config.InitConfiguration()
	database.InitDatabase(conf)
	db := database.DB

	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	authService := auth.NewService()
	userController := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")

	// user
	api.POST("/users", userController.UserRegister)
	api.POST("/sessions", userController.Login)

	router.Run()

}
