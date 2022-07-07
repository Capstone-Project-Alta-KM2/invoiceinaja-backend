package main

import (
	"invoiceinaja/auth"
	"invoiceinaja/config"
	"invoiceinaja/database"
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/payment"

	"github.com/gin-gonic/gin"

	"invoiceinaja/domain/invoice"
	"invoiceinaja/domain/user"
	"invoiceinaja/handler"
	"invoiceinaja/routes"
)

func main() {

	conf := config.InitConfiguration()
	database.InitDatabase(conf)
	db := database.DB

	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)

	clientRepo := client.NewClientRepository(db)
	clientService := client.NewUserService(clientRepo)
	clientHandler := handler.NewClientHandler(clientService, userService, authService)

	paymentService := payment.NewService()
	invoiceRepo := invoice.NewInvoiceRepository(db)
	invoiceService := invoice.NewUserService(invoiceRepo, paymentService)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService, clientService, authService)

	dashboardHandler := handler.NewDashboardHandler(invoiceService, clientService, authService)

	router := gin.Default()

	router.Use(auth.CORSMiddleware())

	// routes.APIRoutes(router, userHandler, clientHandler, authService, userService)
	routes.APIRoutes(
		router,
		userHandler,
		clientHandler,
		invoiceHandler,
		dashboardHandler,
		authService,
		userService,
	)

	router.Run()

}
