package database

import (
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/invoice"
	"invoiceinaja/domain/user"
)

type Model struct {
	Model interface{}
}

func RegisterModel() []Model {
	return []Model{
		{Model: user.User{}},
		{Model: client.Client{}},
		{Model: invoice.Invoice{}},
		{Model: invoice.DetailInvoice{}},
	}
}
