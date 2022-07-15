package invoice

import (
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/invoice"
	"invoiceinaja/domain/payment"
	"invoiceinaja/domain/user"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var repoInvoice = &RepositoryMock{Mock: mock.Mock{}}
var paymentService payment.Service
var serviceInvoice = invoice.NewUserService(repoInvoice, paymentService)

func TestAddInvoice(t *testing.T) {
	inp := invoice.InputAddInvoice{
		Invoice: invoice.InvoiceData{
			ClientID:    1,
			TotalAmount: 4000,
			InvoiceDate: "2022-07-12",
			InvoiceDue:  "2022-07-12",
		},
		DetailInvoice: []invoice.DetailInvoiceData{
			{
				ItemName: "item 1",
				Price:    1000,
				Quantity: 2,
			},
			{
				ItemName: "item 2",
				Price:    2000,
				Quantity: 1,
			},
		},
	}

	outSave := invoice.Invoice{
		ID:          1,
		ClientID:    1,
		Client:      client.Client{ID: 1},
		TotalAmount: 4000,
		PaymentURL:  "",
		Status:      "",
		Items: []invoice.DetailInvoice{
			{
				ID:        1,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 1",
				Price:     1000,
				Quantity:  2,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			{
				ID:        2,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 2",
				Price:     2000,
				Quantity:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		InvoiceDate: "2022-07-12",
		InvoiceDue:  "2022-07-12",
	}

	repoInvoice.Mock.On("Save", mock.AnythingOfType("invoice.Invoice")).Return(outSave)
	res, err := serviceInvoice.AddInvoice(inp)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, outSave.ID, res.ID)
}

func TestCheckDate(t *testing.T) {
	invoices := []invoice.Invoice{
		{
			ID:          1,
			ClientID:    1,
			Client:      client.Client{ID: 1},
			TotalAmount: 4000,
			PaymentURL:  "",
			Status:      "UNPAID",
			Items: []invoice.DetailInvoice{
				{
					ID:        1,
					InvoiceID: 1,
					Invoice:   invoice.Invoice{},
					ItemName:  "item 1",
					Price:     1000,
					Quantity:  2,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				{
					ID:        2,
					InvoiceID: 1,
					Invoice:   invoice.Invoice{},
					ItemName:  "item 2",
					Price:     2000,
					Quantity:  1,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},
			InvoiceDate: "2022-07-12",
			InvoiceDue:  "2022-07-18",
		},
		{
			ID:          2,
			ClientID:    2,
			Client:      client.Client{ID: 1},
			TotalAmount: 4000,
			PaymentURL:  "",
			Status:      "OVERDUE",
			Items: []invoice.DetailInvoice{
				{
					ID:        3,
					InvoiceID: 2,
					Invoice:   invoice.Invoice{},
					ItemName:  "item 1",
					Price:     1000,
					Quantity:  2,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				{
					ID:        4,
					InvoiceID: 2,
					Invoice:   invoice.Invoice{},
					ItemName:  "item 2",
					Price:     2000,
					Quantity:  1,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},
			InvoiceDate: "2022-07-12",
			InvoiceDue:  "2022-07-12",
		},
	}

	repoInvoice.Mock.On("UpdateInvoice", mock.AnythingOfType("invoice.Invoice")).Return(invoices[1])

	res, err := serviceInvoice.CheckDate(invoices)

	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestDeleteDetailInvoice(t *testing.T) {
	inp := invoice.DetailInvoice{
		ID:        1,
		InvoiceID: 1,
		Invoice:   invoice.Invoice{},
		ItemName:  "item 2",
		Price:     2000,
		Quantity:  1,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	repoInvoice.Mock.On("DeleteDetailInvoice", inp).Return(inp)
	res, err := serviceInvoice.DeleteDetailInvoice(inp)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, inp.ID, res.ID)
}

func TestDeleteInvoice(t *testing.T) {
	data := invoice.Invoice{
		ID:          1,
		ClientID:    1,
		Client:      client.Client{ID: 1},
		TotalAmount: 4000,
		PaymentURL:  "",
		Status:      "",
		Items: []invoice.DetailInvoice{
			{
				ID:        1,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 1",
				Price:     1000,
				Quantity:  2,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			{
				ID:        2,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 2",
				Price:     2000,
				Quantity:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		InvoiceDate: "2022-07-12",
		InvoiceDue:  "2022-07-12",
	}

	repoInvoice.Mock.On("FindByID", 1).Return(data)
	repoInvoice.Mock.On("DeleteInvoice", data).Return(data)
	res, err := serviceInvoice.DeleteInvoice(1)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, data.ID, res.ID)
}

func TestGetByID(t *testing.T) {
	data := invoice.Invoice{
		ID:          1,
		ClientID:    1,
		Client:      client.Client{ID: 1},
		TotalAmount: 4000,
		PaymentURL:  "",
		Status:      "",
		Items: []invoice.DetailInvoice{
			{
				ID:        1,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 1",
				Price:     1000,
				Quantity:  2,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			{
				ID:        2,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 2",
				Price:     2000,
				Quantity:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		InvoiceDate: "2022-07-12",
		InvoiceDue:  "2022-07-12",
	}

	repoInvoice.Mock.On("FindByID", 1).Return(data)
	res, err := serviceInvoice.GetByID(1)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, data.ID, res.ID)
}

func TestGetDetailByID(t *testing.T) {
	data := invoice.Invoice{
		ID:          1,
		ClientID:    1,
		Client:      client.Client{ID: 1},
		TotalAmount: 4000,
		PaymentURL:  "",
		Status:      "",
		Items: []invoice.DetailInvoice{
			{
				ID:        1,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 1",
				Price:     1000,
				Quantity:  2,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			{
				ID:        2,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 2",
				Price:     2000,
				Quantity:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		InvoiceDate: "2022-07-12",
		InvoiceDue:  "2022-07-12",
	}

	repoInvoice.Mock.On("FindAllDetail", 1).Return(data.Items)
	res, err := serviceInvoice.GetDetailByID(1)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, data.Items[1].ID, res[1].ID)
}

func TestGetInvoices(t *testing.T) {
	invoices := []invoice.Invoice{
		{
			ID:          1,
			ClientID:    1,
			Client:      client.Client{ID: 1},
			TotalAmount: 4000,
			PaymentURL:  "",
			Status:      "UNPAID",
			Items: []invoice.DetailInvoice{
				{
					ID:        1,
					InvoiceID: 1,
					Invoice:   invoice.Invoice{},
					ItemName:  "item 1",
					Price:     1000,
					Quantity:  2,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				{
					ID:        1,
					InvoiceID: 1,
					Invoice:   invoice.Invoice{},
					ItemName:  "item 2",
					Price:     2000,
					Quantity:  1,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},
			InvoiceDate: "2022-07-12",
			InvoiceDue:  "2022-07-18",
		},
		{
			ID:          2,
			ClientID:    2,
			Client:      client.Client{ID: 1},
			TotalAmount: 4000,
			PaymentURL:  "",
			Status:      "OVERDUE",
			Items: []invoice.DetailInvoice{
				{
					ID:        3,
					InvoiceID: 2,
					Invoice:   invoice.Invoice{},
					ItemName:  "item 1",
					Price:     1000,
					Quantity:  2,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				{
					ID:        4,
					InvoiceID: 2,
					Invoice:   invoice.Invoice{},
					ItemName:  "item 2",
					Price:     2000,
					Quantity:  1,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},
			InvoiceDate: "2022-07-12",
			InvoiceDue:  "2022-07-12",
		},
	}

	repoInvoice.Mock.On("FindAll", 1).Return(invoices)
	res, err := serviceInvoice.GetInvoices(1)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, invoices[1].ID, res[1].ID)
}

func TestGetSumStatus(t *testing.T) {
	result := map[string]int{"paid": 0, "unpaid": 0}

	repoInvoice.Mock.On("FindPaid", 1).Return(result)
	res, err := serviceInvoice.GetSumStatus(1)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, result["paid"], res["paid"])
}

func TestProcessPayment(t *testing.T) {
	inp := payment.InputTransactionNotif{
		TransactionStatus: "capture",
		OrderID:           "1",
		PaymentType:       "credit_card",
		FraudStatus:       "accept",
	}

	inv := invoice.Invoice{
		ID:          2,
		ClientID:    1,
		Client:      client.Client{ID: 1},
		TotalAmount: 4000,
		PaymentURL:  "",
		Status:      "UNPAID",
		Items: []invoice.DetailInvoice{
			{
				ID:        1,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 1",
				Price:     1000,
				Quantity:  2,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			{
				ID:        2,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 2",
				Price:     2000,
				Quantity:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		InvoiceDate: "2022-07-12",
		InvoiceDue:  "2022-07-18",
	}

	invPaid := invoice.Invoice{
		ID:          2,
		ClientID:    1,
		Client:      client.Client{ID: 1},
		TotalAmount: 4000,
		PaymentURL:  "",
		Status:      "PAID",
		Items: []invoice.DetailInvoice{
			{
				ID:        1,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 1",
				Price:     1000,
				Quantity:  2,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			{
				ID:        2,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 2",
				Price:     2000,
				Quantity:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		InvoiceDate: "2022-07-12",
		InvoiceDue:  "2022-07-18",
	}

	repoInvoice.Mock.On("FindByID", 2).Return(inv)
	repoInvoice.Mock.On("UpdateInvoice", mock.AnythingOfType("invoice.Invoice")).Return(invPaid)
	res, err := serviceInvoice.ProcessPayment(inp)
	res.Status = "PAID"

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.Status, invPaid.Status)
}

func TestSaveDetail(t *testing.T) {
	inp := invoice.InputAddInvoice{
		Invoice: invoice.InvoiceData{
			ClientID:    1,
			TotalAmount: 4000,
			InvoiceDate: "2022-07-12",
			InvoiceDue:  "2022-07-12",
		},
		DetailInvoice: []invoice.DetailInvoiceData{
			{
				ItemName: "item 1",
				Price:    1000,
				Quantity: 2,
			},
		},
	}

	asd := []invoice.DetailInvoice{
		{
			ID:        0,
			InvoiceID: 1,
			Invoice:   invoice.Invoice{},
			ItemName:  "item 1",
			Price:     1000,
			Quantity:  2,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}

	repoInvoice.Mock.On("SaveDetail", asd).Return(asd)
	res, err := serviceInvoice.SaveDetail(1, inp)

	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestSendMailInvoice(t *testing.T) {
	inv := invoice.Invoice{
		ID:          1,
		ClientID:    1,
		Client:      client.Client{ID: 1},
		TotalAmount: 4000,
		PaymentURL:  "",
		Status:      "",
		Items: []invoice.DetailInvoice{
			{
				ID:        1,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 1",
				Price:     1000,
				Quantity:  2,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			{
				ID:        2,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 2",
				Price:     2000,
				Quantity:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		InvoiceDate: "2022-07-12",
		InvoiceDue:  "2022-07-12",
	}

	usr := user.User{
		ID:           1,
		Fullname:     "fahmi",
		Email:        "fahmi@gmail.com",
		NoTlpn:       "08521",
		Password:     "password",
		BusinessName: "ojol",
		Role:         "user",
		Avatar:       "image/avatar.png",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}

	cln := client.Client{
		ID:        1,
		Fullname:  "client2",
		Email:     "client2@gmail.com",
		Address:   "bandung",
		City:      "bandung",
		ZipCode:   "41523",
		Company:   "toko sembako",
		UserID:    1,
		User:      user.User{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}

	repoInvoice.Mock.On("FindByID", 1).Return(inv)
	res, err := serviceInvoice.SendMailInvoice(1, usr, cln)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res, inv)
}

func TestUpdateInvoice(t *testing.T) {
	inv := invoice.Invoice{
		ID:          2,
		ClientID:    2,
		Client:      client.Client{ID: 1},
		TotalAmount: 4000,
		PaymentURL:  "",
		Status:      "UNPAID",
		Items: []invoice.DetailInvoice{
			{
				ID:        1,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 1",
				Price:     1000,
				Quantity:  2,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			{
				ID:        2,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 2",
				Price:     2000,
				Quantity:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		InvoiceDate: "2022-07-12",
		InvoiceDue:  "2022-07-12",
	}
	invP := invoice.Invoice{
		ID:          2,
		ClientID:    2,
		Client:      client.Client{ID: 1},
		TotalAmount: 4000,
		PaymentURL:  "url",
		Status:      "PAID",
		Items: []invoice.DetailInvoice{
			{
				ID:        1,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 1",
				Price:     1000,
				Quantity:  2,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			{
				ID:        2,
				InvoiceID: 1,
				Invoice:   invoice.Invoice{},
				ItemName:  "item 2",
				Price:     2000,
				Quantity:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		InvoiceDate: "2022-07-12",
		InvoiceDue:  "2022-07-12",
	}

	repoInvoice.Mock.On("UpdateInvoice", inv).Return(invP)
	res, err := serviceInvoice.UpdateInvoice(inv, "url")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.PaymentURL, inv.PaymentURL)
}
