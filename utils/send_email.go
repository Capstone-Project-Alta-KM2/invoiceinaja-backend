package utils

import (
	"log"

	"github.com/mailjet/mailjet-apiv3-go"
)

func SendMail(destination, newPassword string) *mailjet.ResultsV31 {
	mailjetClient := mailjet.NewMailjetClient("5f4b8dba26ef85efb6dce6410157bbe9", "efd86ba2c3502512da935ad19de63869")
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "kokolopo321@gmail.com",
				Name:  "InvoiceinAja",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: destination,
					Name:  destination,
				},
			},
			Subject:  "Hay Aden ini InvoiceinAja.",
			TextPart: "hay ini email make mailjet",
			HTMLPart: "<h3>Dear passenger 1, welcome to <a href='https://www.mailjet.com/'>Mailjet</a>!</h3><br />May the delivery force be with you!" + newPassword,
			CustomID: "AppGettingStartedTest",
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Data: %+v\n", res)
	return res
}
