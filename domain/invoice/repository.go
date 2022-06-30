package invoice

import (
	"gorm.io/gorm"
)

type IRepository interface {
	Save(invoice Invoice) (Invoice, error)
	SaveDetail(detail []DetailInvoice) ([]DetailInvoice, error)
	FindAll(userID int) ([]Invoice, error)
	// FindById(id int) (Invoice, error)
	// FindByEmail(email string) (Invoice, error)
	// Update(invoice Invoice) (Invoice, error)
	// Delete(invoice Invoice) (Invoice, error)
}

type repository struct {
	DB *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(invoice Invoice) (Invoice, error) {
	err := r.DB.Create(&invoice).Error
	if err != nil {
		return invoice, err
	}

	return invoice, nil
}

func (r *repository) SaveDetail(detail []DetailInvoice) ([]DetailInvoice, error) {
	err := r.DB.Create(&detail).Error
	if err != nil {
		return detail, err
	}

	return detail, nil
}

func (r *repository) FindAll(userID int) ([]Invoice, error) {
	var invoices []Invoice
	var invoicesByUser []Invoice
	err := r.DB.Preload("Client").Preload("Items", "detail_invoices.invoice_id").Order("id desc").Find(&invoices).Error
	if err != nil {
		return invoices, err
	}

	for _, v := range invoices {
		if v.Client.UserID == userID {
			invoicesByUser = append(invoicesByUser, v)
		}
	}

	return invoicesByUser, nil
}
