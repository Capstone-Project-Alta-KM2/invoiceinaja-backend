package invoice

import (
	"time"

	"invoiceinaja/domain/client"
)

type Invoice struct {
	ID          int           `gorm:"primary_key;auto_increment;not_null"`
	ClientID    int           `gorm:"type:int(25);not null"`
	Client      client.Client `gorm:"foreignKey:ClientID;not null"`
	TotalAmount int           `gorm:"type:int(100);not null"`
	PaymentURL  string        `gorm:"type:varchar(100);not null"`
	Status      string        `gorm:"type:varchar(100);not null"`
	Items       []DetailInvoice
	InvoiceDate string
	InvoiceDue  string
}

type DetailInvoice struct {
	ID        int     `gorm:"primary_key;auto_increment;not_null"`
	InvoiceID int     `gorm:"type:int(25);not null"`
	Invoice   Invoice `gorm:"foreignKey:InvoiceID;not null"`
	ItemName  string  `gorm:"type:varchar(100);not null"`
	Price     int     `gorm:"type:int(100);not null"`
	Quantity  int     `gorm:"type:int(25);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DataCSV struct {
	Name        string
	Email       string
	Item        string
	Price       string
	Quantity    string
	InvoiceDate string
	InvoiceDue  string
}

type InvoiceCSV struct {
	ID          int
	Name        string
	Email       string
	Items       []Item
	InvoiceDate string
	InvoiceDue  string
}

type Item struct {
	ItemName string
	Price    int
	Quantity int
}
