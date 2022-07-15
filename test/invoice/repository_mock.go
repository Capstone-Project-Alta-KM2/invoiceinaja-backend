package invoice

import (
	"errors"
	"invoiceinaja/domain/invoice"

	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	Mock mock.Mock
}

func NewRepositoryMock(mock mock.Mock) *RepositoryMock {
	return &RepositoryMock{mock}
}

// Save(invoice Invoice) (Invoice, error)
// 	SaveDetail(detail []DetailInvoice) ([]DetailInvoice, error)
// 	FindAll(userID int) ([]Invoice, error)
// 	FindByID(id int) (Invoice, error)
// 	FindAllDetail(invoiceID int) ([]DetailInvoice, error)
// 	FindPaid(userID int) (map[string]int, error)
// 	DeleteInvoice(invoice Invoice) (Invoice, error)
// 	DeleteDetailInvoice(detailInvoice DetailInvoice) (DetailInvoice, error)
// 	UpdateInvoice(invoice Invoice) (Invoice, error)

func (r *RepositoryMock) Save(params invoice.Invoice) (invoice.Invoice, error) {
	argument := r.Mock.Called(params)
	if argument.Get(0) == nil {
		return params, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).(invoice.Invoice)
		return client, nil
	}
}

func (r *RepositoryMock) SaveDetail(params []invoice.DetailInvoice) ([]invoice.DetailInvoice, error) {
	argument := r.Mock.Called(params)
	if argument.Get(0) == nil {
		return params, errors.New("ada yang salah")
	} else {
		detail := argument.Get(0).([]invoice.DetailInvoice)
		return detail, nil
	}
}

func (r *RepositoryMock) FindAll(userID int) ([]invoice.Invoice, error) {
	argument := r.Mock.Called(userID)
	if argument.Get(0) == nil {
		return []invoice.Invoice{}, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).([]invoice.Invoice)
		return client, nil
	}
}

func (r *RepositoryMock) FindByID(id int) (invoice.Invoice, error) {
	argument := r.Mock.Called(id)
	if argument.Get(0) == nil {
		return invoice.Invoice{}, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).(invoice.Invoice)
		return client, nil
	}
}

func (r *RepositoryMock) FindAllDetail(invoiceID int) ([]invoice.DetailInvoice, error) {
	argument := r.Mock.Called(invoiceID)
	if argument.Get(0) == nil {
		return []invoice.DetailInvoice{}, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).([]invoice.DetailInvoice)
		return client, nil
	}
}

func (r *RepositoryMock) FindPaid(userID int) (map[string]int, error) {
	argument := r.Mock.Called(userID)
	if argument.Get(0) == nil {
		return map[string]int{}, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).(map[string]int)
		return client, nil
	}
}

func (r *RepositoryMock) DeleteInvoice(params invoice.Invoice) (invoice.Invoice, error) {
	argument := r.Mock.Called(params)
	if argument.Get(0) == nil {
		return invoice.Invoice{}, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).(invoice.Invoice)
		return client, nil
	}
}

func (r *RepositoryMock) DeleteDetailInvoice(detailInvoice invoice.DetailInvoice) (invoice.DetailInvoice, error) {
	argument := r.Mock.Called(detailInvoice)
	if argument.Get(0) == nil {
		return invoice.DetailInvoice{}, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).(invoice.DetailInvoice)
		return client, nil
	}
}

func (r *RepositoryMock) UpdateInvoice(params invoice.Invoice) (invoice.Invoice, error) {
	argument := r.Mock.Called(params)
	if argument.Get(0) == nil {
		return invoice.Invoice{}, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).(invoice.Invoice)
		return client, nil
	}
}
