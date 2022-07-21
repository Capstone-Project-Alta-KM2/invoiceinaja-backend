package invoice

import (
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/user"
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

