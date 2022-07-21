package routes

import (
	"net/http"

	"invoiceinaja/auth"
	"invoiceinaja/domain/user"
	"invoiceinaja/handler"

	"github.com/gin-gonic/gin"
)

func APIRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	clientHandler *handler.ClientHandler,
	invoiceHandler *handler.InvoiceHandler,
	dashboardHandler *handler.DashboardHandler,
	authService auth.Service,
	userService user.IService) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Wellcome to InvoiceinAja !!!!!",
		})
	})

	router.Static("/images", "./images")

	api := router.Group("/api/v1")

	// user
	api.POST("/users", userHandler.UserRegister)
	api.GET("/resend_otp", auth.AuthMiddleware(authService, userService), userHandler.ResendOTP)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/sessions", userHandler.Login)
	api.POST("/avatars", auth.AuthMiddleware(authService, userService), userHandler.UploadAvatar)
	api.PUT("/users", auth.AuthMiddleware(authService, userService), userHandler.UpdateUser)
	api.POST("/change_passwords", auth.AuthMiddleware(authService, userService), userHandler.ChangePassword)
	api.POST("/reset_passwords", userHandler.ResetPassword)

	// client
	api.POST("/clients", auth.AuthMiddleware(authService, userService), clientHandler.AddClient)
	api.POST("/clients_by_csv", auth.AuthMiddleware(authService, userService), clientHandler.AddClientsByCSV)
	api.GET("/clients", auth.AuthMiddleware(authService, userService), clientHandler.GetClients)
	api.PUT("/clients/:id", auth.AuthMiddleware(authService, userService), clientHandler.UpdateClient)
	api.DELETE("/clients/:id", auth.AuthMiddleware(authService, userService), clientHandler.DeleteClient)
	api.GET("/clients/invoices/:id", invoiceHandler.GetClientInvoice)

	// invoice
	api.POST("/invoices", auth.AuthMiddleware(authService, userService), invoiceHandler.AddInvoice)
	api.POST("/invoices_by_csv", auth.AuthMiddleware(authService, userService), invoiceHandler.GenerateByCSV)
	api.GET("/invoices", auth.AuthMiddleware(authService, userService), invoiceHandler.GetInvoices)
	api.GET("/invoices/:id", auth.AuthMiddleware(authService, userService), invoiceHandler.GetInvoicesByID)
	api.DELETE("/invoices/:id", auth.AuthMiddleware(authService, userService), invoiceHandler.DeleteInvoice)
	api.POST("/invoice_payments", invoiceHandler.InvoicePay)
	api.POST("/invoice_payments/notification", invoiceHandler.GetNotification)
	// api.POST("/reminders", auth.AuthMiddleware(authService,userService), invoiceHandler.SendReminder)

	// dashboard
	api.GET("/overall", auth.AuthMiddleware(authService, userService), dashboardHandler.GetDataOverall)
	api.GET("/graphics", auth.AuthMiddleware(authService, userService), dashboardHandler.GetDataGraphic)
}
