package invoice

import (
	"log"

	"github.com/mailjet/mailjet-apiv3-go"
)

type InputAddInvoice struct {
	Invoice       InvoiceData         `json:"invoice" binding:"required"`
	DetailInvoice []DetailInvoiceData `json:"detail_invoice" binding:"required"`
}

type InvoiceData struct {
	ClientID    int    `json:"client_id" binding:"required"`
	TotalAmount int    `json:"total_amount" binding:"required"`
	InvoiceDate string `json:"invoice_date" binding:"required"`
	InvoiceDue  string `json:"invoice_due" binding:"required"`
}

type DetailInvoiceData struct {
	ItemName string `json:"item_name" binding:"required"`
	Price    int    `json:"price" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

func SendMailInvoice(recaiver string, data Invoice) *mailjet.ResultsV31 {
	mailjetClient := mailjet.NewMailjetClient("5f4b8dba26ef85efb6dce6410157bbe9", "efd86ba2c3502512da935ad19de63869")
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "kokolopo321@gmail.com",
				Name:  "InvoiceinAja",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: recaiver,
					Name:  recaiver,
				},
			},
			TemplateID:       4036084,
			TemplateLanguage: true,
			Subject:          "Your Invoice",
			Variables: map[string]interface{}{
				"status":        data.Status,
				"invoiceNo":     data.ID,
				"userName":      "nama user",
				"userCompany":   "user comp",
				"userEmail":     "user email",
				"client":        data.Client,
				"clientEmail":   "client email",
				"clientCompany": "client comp",
			},
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Data: %+v", res)
	return res
}
