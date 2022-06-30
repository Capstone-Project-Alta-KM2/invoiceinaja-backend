package invoice

import (
	"errors"
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/user"
)

type IService interface {
	AddInvoice(input InputAddInvoice) (Invoice, error)
	SaveDetail(invoiceID int, input InputAddInvoice) ([]DetailInvoice, error)
	GetInvoices(userID int) ([]Invoice, error)
	SendMailInvoice(invoiceID int, user user.User, client client.Client) (Invoice, error)
	// GetAll(clientID int) ([]Client, error)
	// GetByID(clientID int) (Client, error)
	// DeleteClient(clientID int) (Client, error)
}

type service struct {
	repository IRepository
}

func NewUserService(repository IRepository) *service {
	return &service{repository}
}

func (s *service) AddInvoice(input InputAddInvoice) (Invoice, error) {
	var invoiceData Invoice

	invoiceData.ClientID = input.Invoice.ClientID
	invoiceData.TotalAmount = input.Invoice.TotalAmount
	invoiceData.InvoiceDate = input.Invoice.InvoiceDate
	invoiceData.InvoiceDue = input.Invoice.InvoiceDue
	invoiceData.Status = "UNPAID"

	invoice, errInvoice := s.repository.Save(invoiceData)
	if errInvoice != nil {
		return invoiceData, errInvoice
	}

	return invoice, errInvoice
}

func (s *service) SaveDetail(invoiceID int, input InputAddInvoice) ([]DetailInvoice, error) {
	var detail []DetailInvoice

	for _, v := range input.DetailInvoice {
		detail = append(detail, DetailInvoice{InvoiceID: invoiceID, ItemName: v.ItemName, Price: v.Price, Quantity: v.Quantity})
	}

	//save data yang sudah dimapping kedalam struct DetailOrder
	newDetail, err := s.repository.SaveDetail(detail)
	if err != nil {
		return newDetail, err
	}

	return newDetail, nil
}

func (s *service) GetInvoices(userID int) ([]Invoice, error) {
	clients, err := s.repository.FindAll(userID)
	if err != nil {
		return clients, err
	}

	return clients, nil
}

func (s *service) SendMailInvoice(invoiceID int, user user.User, client client.Client) (Invoice, error) {
	var data SendEmailData

	invoice, err := s.repository.FindByID(invoiceID)
	if err != nil {
		return invoice, err
	}

	if invoice.ID == 0 {
		return invoice, errors.New("no invoice found")
	}

	data.Invoice.ID = invoice.ID
	data.Invoice.Status = invoice.Status
	data.Invoice.TotalAmount = invoice.TotalAmount
	data.Invoice.InvoiceDate = invoice.InvoiceDate
	data.Invoice.InvoiceDue = invoice.InvoiceDue
	data.User.Fullname = user.Fullname
	data.User.BusinessName = user.BusinessName
	data.User.Email = user.Email
	data.Client.Fullname = client.Fullname
	data.Client.Email = client.Email
	data.Client.Address = client.Address
	data.Client.City = client.City
	data.Client.ZipCode = client.ZipCode
	data.Client.Company = client.Company
	data.Invoice.Items = append(data.Invoice.Items, invoice.Items...)

	SendMailInvoice(data.Client.Email, data)

	return invoice, nil
}
