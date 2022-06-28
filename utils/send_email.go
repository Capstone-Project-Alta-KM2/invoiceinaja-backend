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
			TemplateID:       4035911,
			TemplateLanguage: true,
			Subject:          "Reset Password InvoiceinAja",
			Variables:        map[string]interface{}{"new_password": newPassword},
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
