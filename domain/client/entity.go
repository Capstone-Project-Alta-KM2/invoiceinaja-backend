package client

import (
	"time"

	"invoiceinaja/domain/user"

	"gorm.io/gorm"
)

type Client struct {
	ID        int       `gorm:"primary_key;auto_increment;not_null"`
	Fullname  string    `gorm:"type:varchar(50);not null"`
	Email     string    `gorm:"type:varchar(100);not null"`
	Address   string    `gorm:"type:longtext;not null"`
	City      string    `gorm:"type:varchar(50);not null"`
	ZipCode   string    `gorm:"type:varchar(50);not null"`
	Company   string    `gorm:"type:varchar(50);not null"`
	UserID    int       `gorm:"type:int(25);not null"`
	User      user.User `gorm:"foreignKey:UserID;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
