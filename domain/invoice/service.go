package invoice

import (
	"errors"
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/payment"
	"invoiceinaja/domain/user"
	"invoiceinaja/helper"
	"invoiceinaja/utils"
	"strconv"
	"time"
)

type IService interface {
	AddInvoice(input InputAddInvoice) (Invoice, error)
	SaveDetail(invoiceID int, input InputAddInvoice) ([]DetailInvoice, error)
	GetInvoices(userID int) ([]Invoice, error)
	SendMailInvoice(invoiceID int, user user.User, client client.Client) (Invoice, error)
	GetSumStatus(userID int) (map[string]int, error)
	GetByID(invoiceID int) (Invoice, error)
	GetDetailByID(invoiceID int) ([]DetailInvoice, error)
	DeleteInvoice(invoiceID int) (Invoice, error)
	DeleteDetailInvoice(detailInvoice DetailInvoice) (DetailInvoice, error)
	PayInvoice(input payment.InputCreateTansaction, client client.Client) (string, error)
	UpdateInvoice(invoice Invoice, urlPayment string) (Invoice, error)
	ProcessPayment(input payment.InputTransactionNotif) (Invoice, error)
}

type service struct {
	repository     IRepository
	paymentService payment.Service
}

func NewUserService(repository IRepository, paymentService payment.Service) *service {
	return &service{repository, paymentService}
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
	invoices, err := s.repository.FindAll(userID)
	if err != nil {
		return invoices, err
	}

	return invoices, nil
}

func (s *service) GetByID(invoiceID int) (Invoice, error) {
	invoice, err := s.repository.FindByID(invoiceID)
	if err != nil {
		return invoice, err
	}

	return invoice, nil
}

func (s *service) GetDetailByID(invoiceID int) ([]DetailInvoice, error) {
	details, err := s.repository.FindAllDetail(invoiceID)
	if err != nil {
		return details, err
	}

	return details, nil
}

func (s *service) DeleteInvoice(invoiceID int) (Invoice, error) {
	invoice, err := s.repository.FindByID(invoiceID)
	if err != nil {
		return invoice, err
	}

	deleteInvoice, errDel := s.repository.DeleteInvoice(invoice)
	if errDel != nil {
		return deleteInvoice, errDel
	}

	return deleteInvoice, nil
}

func (s *service) DeleteDetailInvoice(detailInvoice DetailInvoice) (DetailInvoice, error) {
	deleteInvoice, errDel := s.repository.DeleteDetailInvoice(detailInvoice)
	if errDel != nil {
		return deleteInvoice, errDel
	}

	return deleteInvoice, nil
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

func (s *service) GetSumStatus(userID int) (map[string]int, error) {
	total, err := s.repository.FindPaid(userID)
	if err != nil {
		return total, err
	}

	return total, nil
}

func (s *service) UpdateInvoice(invoice Invoice, urlPayment string) (Invoice, error) {
	invoice.PaymentURL = urlPayment

	updatedItem, errUpdate := s.repository.UpdateInvoice(invoice)
	if errUpdate != nil {
		return updatedItem, errUpdate
	}

	return updatedItem, nil
}

func (s *service) PayInvoice(input payment.InputCreateTansaction, client client.Client) (string, error) {
	paymentTransaction := payment.Transaction{
		ID:     input.InvoiceID,
		Amount: input.TotalAmount,
	}

	urlPayment, err := s.paymentService.GetPaymentUrl(paymentTransaction, client)
	if err != nil {
		return urlPayment, err
	}

	return urlPayment, nil
}

func (s *service) ProcessPayment(input payment.InputTransactionNotif) (Invoice, error) {
	transaction_id, _ := strconv.Atoi(input.OrderID)

	invoice, err := s.repository.FindByID(transaction_id)
	if err != nil {
		return invoice, err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		invoice.Status = "PAID"
	} else if input.TransactionStatus == "settlement" {
		invoice.Status = "PAID"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		invoice.Status = "CANCELLED"
	}

	updatedInvoice, err := s.repository.UpdateInvoice(invoice)
	if err != nil {
		return updatedInvoice, err
	}

	return updatedInvoice, nil
}

func Mapping(lines [][]string) []InvoiceCSV {
	var inv []InvoiceCSV
	var email []string // nampung email

	// Loop through lines & turn into object
	for i, line := range lines {
		if i == 0 {
			continue
		}

		// mapping data
		data := DataCSV{
			Name:        line[0],
			Email:       line[1],
			Item:        line[2],
			Price:       line[3],
			Quantity:    line[4],
			InvoiceDate: line[5],
			InvoiceDue:  line[6],
		}

		// if email sudah ada pada slice email
		// email tidak akan diappend kedalam slice email
		isExist := utils.Contains(email, data.Email)
		switch isExist {
		case true:
			var a int
			for idx, v := range email {
				if data.Email == v {
					a = idx
				}
			}
			inv[a].Items = append(inv[a].Items, Item{ItemName: data.Item, Price: ConvStrToInt(data.Price), Quantity: ConvStrToInt(data.Quantity)})
		case false:
			email = append(email, data.Email)
			inv = append(inv, InvoiceCSV{
				Name:        data.Name,
				Email:       data.Email,
				Items:       []Item{{ItemName: data.Item, Price: ConvStrToInt(data.Price), Quantity: ConvStrToInt(data.Quantity)}},
				InvoiceDate: data.InvoiceDate,
				InvoiceDue:  data.InvoiceDue,
			})
		}
	}
	return inv
}

func ConvStrToInt(data string) int {
	i, _ := strconv.Atoi(data)
	return i
}

func GraphicInvoice(invoices []Invoice) helper.Month {
	var s helper.Month

	for _, v := range invoices {
		x := 0
		switch {
		case helper.ConvStingDate(v.InvoiceDue).Month().String() == time.January.String():
			if v.Status == "PAID" {
				x += v.TotalAmount
				s.Jan.Paid += x
			}
			if v.Status == "UNPAID" || v.Status == "OVERDUE" {
				x += v.TotalAmount
				s.Jan.Unpaid += x
			}
		case helper.ConvStingDate(v.InvoiceDue).Month().String() == time.February.String():
			if v.Status == "PAID" {
				x += v.TotalAmount
				s.Feb.Paid += x
			}
			if v.Status == "UNPAID" || v.Status == "OVERDUE" {
				x += v.TotalAmount
				s.Feb.Unpaid += x
			}
		case helper.ConvStingDate(v.InvoiceDue).Month().String() == time.March.String():
			if v.Status == "PAID" {
				x += v.TotalAmount
				s.Mar.Paid += x
			}
			if v.Status == "UNPAID" || v.Status == "OVERDUE" {
				x += v.TotalAmount
				s.Mar.Unpaid += x
			}
		case helper.ConvStingDate(v.InvoiceDue).Month().String() == time.April.String():
			if v.Status == "PAID" {
				x += v.TotalAmount
				s.Apr.Paid += x
			}
			if v.Status == "UNPAID" || v.Status == "OVERDUE" {
				x += v.TotalAmount
				s.Apr.Unpaid += x
			}
		case helper.ConvStingDate(v.InvoiceDue).Month().String() == time.May.String():
			if v.Status == "PAID" {
				x += v.TotalAmount
				s.May.Paid += x
			}
			if v.Status == "UNPAID" || v.Status == "OVERDUE" {
				x += v.TotalAmount
				s.May.Unpaid += x
			}
		case helper.ConvStingDate(v.InvoiceDue).Month().String() == time.June.String():
			if v.Status == "PAID" {
				x += v.TotalAmount
				s.Jun.Paid += x
			}
			if v.Status == "UNPAID" || v.Status == "OVERDUE" {
				x += v.TotalAmount
				s.Jun.Unpaid += x
			}
		case helper.ConvStingDate(v.InvoiceDue).Month().String() == time.July.String():
			if v.Status == "PAID" {
				x += v.TotalAmount
				s.Jul.Paid += x
			}
			if v.Status == "UNPAID" || v.Status == "OVERDUE" {
				x += v.TotalAmount
				s.Jul.Unpaid += x
			}
		}
	}
	return s
}
