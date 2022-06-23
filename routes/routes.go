package routes

import (
	"net/http"

	"invoiceinaja/auth"
	"invoiceinaja/domain/user"
	"invoiceinaja/handler"

	"github.com/gin-gonic/gin"
)

func APIRoutes(router *gin.Engine, userHandler *handler.UserHandler, clientHandler *handler.ClientHandler, invoiceHandler *handler.InvoiceHandler, authService auth.Service, userService user.IService) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "wellcome",
		})
	})

	api := router.Group("/api/v1")

	// user
	api.POST("/users", userHandler.UserRegister)
	//api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/sessions", userHandler.Login)
	api.POST("/avatars", auth.AuthMiddleware(authService, userService), userHandler.UploadAvatar)
	api.PUT("/users", auth.AuthMiddleware(authService, userService), userHandler.UpdateUser)
	api.POST("/reset_passwords", userHandler.ResetPassword)

	// client
	api.POST("/clients", auth.AuthMiddleware(authService, userService), clientHandler.AddClient)
	api.GET("/clients", auth.AuthMiddleware(authService, userService), clientHandler.GetClients)
	api.PUT("/clients", auth.AuthMiddleware(authService, userService), clientHandler.UpdateClient)
	api.DELETE("/clients/:id", auth.AuthMiddleware(authService, userService), clientHandler.DeleteClient)

	// invoice
	api.POST("/invoices", auth.AuthMiddleware(authService, userService), invoiceHandler.AddInvoice)
	api.GET("/invoices", auth.AuthMiddleware(authService, userService), invoiceHandler.GetInvoices)
}
