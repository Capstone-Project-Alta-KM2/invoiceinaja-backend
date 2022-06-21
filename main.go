package main

import (
	"invoiceinaja/auth"
	"invoiceinaja/config"
	"invoiceinaja/database"
	"invoiceinaja/domain/user"
	"invoiceinaja/handler"
	"invoiceinaja/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.InitConfiguration()
	database.InitDatabase(conf)
	db := database.DB

	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()

	router.Use(auth.CORSMiddleware())

	routes.APIRoutes(router, userHandler, authService, userService)

	router.Run()

}
