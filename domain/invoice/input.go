package invoice

import (
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/user"
	"log"
	"strconv"

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

type SendEmailData struct {
	Invoice Invoice
	User    user.User
	Client  client.Client
}

func SendMailInvoice(recaiver string, data SendEmailData) *mailjet.ResultsV31 {
	mailjetClient := mailjet.NewMailjetClient("5f4b8dba26ef85efb6dce6410157bbe9", "efd86ba2c3502512da935ad19de63869")

	totalAmount := 0
	length := len(data.Invoice.Items)
	var variableHtml = map[string]interface{}{
		"status":         data.Invoice.Status,
		"id_invoice":     data.Invoice.ID,
		"user_name":      data.User.Fullname,
		"date":           data.Invoice.InvoiceDate,
		"due":            data.Invoice.InvoiceDue,
		"user_company":   data.User.BusinessName,
		"user_email":     data.User.Email,
		"client_name":    data.Client.Fullname,
		"client_email":   data.Client.Email,
		"client_address": data.Client.Address,
		"client_city":    data.Client.City,
		"client_zip":     data.Client.ZipCode,
		"client_company": data.Client.Company,
	}
	for i := 0; i < length; i++ {
		variableHtml["item"+strconv.Itoa(i+1)] = data.Invoice.Items[i].ItemName
		variableHtml["price"+strconv.Itoa(i+1)] = data.Invoice.Items[i].Price
		variableHtml["quantity"+strconv.Itoa(i+1)] = data.Invoice.Items[i].Quantity
		variableHtml["total"+strconv.Itoa(i+1)] = data.Invoice.Items[i].Price * data.Invoice.Items[i].Quantity

		totalAmount += data.Invoice.Items[i].Price * data.Invoice.Items[i].Quantity
	}
	variableHtml["total_amount"] = totalAmount

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
			Variables:        variableHtml,
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
