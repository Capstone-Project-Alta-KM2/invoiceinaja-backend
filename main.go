package main

// import (
// 	"invoiceinaja/auth"
// 	"invoiceinaja/config"
// 	"invoiceinaja/database"
// 	"invoiceinaja/domain/user"
// 	"invoiceinaja/handler"
// 	"invoiceinaja/routes"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	conf := config.InitConfiguration()
// 	database.InitDatabase(conf)
// 	db := database.DB

// 	userRepo := user.NewUserRepository(db)
// 	userService := user.NewUserService(userRepo)
// 	authService := auth.NewService()
// 	userHandler := handler.NewUserHandler(userService, authService)

// 	router := gin.Default()

// 	router.Use(auth.CORSMiddleware())

// 	routes.APIRoutes(router, userHandler, authService, userService)

// 	router.Run()

// }

import (
	"fmt"
	"log"

	"github.com/mailjet/mailjet-apiv3-go"
)

func main() {
	mailjetClient := mailjet.NewMailjetClient("5f4b8dba26ef85efb6dce6410157bbe9", "efd86ba2c3502512da935ad19de63869")
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "kokolopo321@gmail.com",
				Name:  "koko",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: "syahrilanas09@gmail.com",
					Name:  "koko",
				},
			},
			Subject:  "Hay Aden ini Kanjut.",
			TextPart: "hay ini email make mailjet",
			HTMLPart: "<h3>Dear passenger 1, welcome to <a href='https://www.mailjet.com/'>Mailjet</a>!</h3><br />May the delivery force be with you!",
			CustomID: "AppGettingStartedTest",
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data: %+v\n", res)
}
