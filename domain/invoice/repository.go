package invoice

import (
	"gorm.io/gorm"
)

type IRepository interface {
	Save(invoice Invoice) (Invoice, error)
	SaveDetail(detail []DetailInvoice) ([]DetailInvoice, error)
	FindAll(userID int) ([]Invoice, error)
	FindByID(id int) (Invoice, error)
	FindAllDetail(invoiceID int) ([]DetailInvoice, error)
	FindPaid(userID int) (map[string]int, error)
	DeleteInvoice(invoice Invoice) (Invoice, error)
	DeleteDetailInvoice(detailInvoice DetailInvoice) (DetailInvoice, error)
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

func (r *repository) FindByID(id int) (Invoice, error) {
	var invoice Invoice
	err := r.DB.Where("id = ?", id).Preload("Client").Preload("Items", "detail_invoices.invoice_id").Order("id desc").Find(&invoice).Error
	if err != nil {
		return invoice, err
	}

	return invoice, nil
}

func (r *repository) FindPaid(userID int) (map[string]int, error) {
	var invoices []Invoice
	var invoicesByUser []Invoice
	result := map[string]int{"paid": 0, "unpaid": 0}
	err := r.DB.Preload("Client").Preload("Items", "detail_invoices.invoice_id").Order("id desc").Find(&invoices).Error
	if err != nil {
		return result, err
	}

	for _, v := range invoices {
		if v.Client.UserID == userID {
			invoicesByUser = append(invoicesByUser, v)
		}
	}

	for _, v := range invoicesByUser {
		if v.Status == "PAID" {
			result["paid"] += v.TotalAmount
		}
		if v.Status == "UNPAID" {
			result["unpaid"] += v.TotalAmount
		}
	}

	return result, nil
}

func (r *repository) FindAllDetail(invoiceID int) ([]DetailInvoice, error) {
	var detail []DetailInvoice
	err := r.DB.Where("invoice_id = ?", invoiceID).Find(&detail).Error
	if err != nil {
		return detail, err
	}

	return detail, nil
}

func (r *repository) DeleteInvoice(invoice Invoice) (Invoice, error) {
	err := r.DB.Delete(&invoice).Error
	if err != nil {
		return invoice, err
	}

	return invoice, nil
}

func (r *repository) DeleteDetailInvoice(detailInvoice DetailInvoice) (DetailInvoice, error) {
	err := r.DB.Delete(&detailInvoice).Error
	if err != nil {
		return detailInvoice, err
	}

	return detailInvoice, nil
}
