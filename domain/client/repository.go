package client

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type IRepository interface {
	Save(client Client) (Client, error)
	FindAll(s string, userID, page, perPage int) ([]Client, int, error)
	FindById(id int) (Client, error)
	FindByEmail(email string, userID int) (Client, error)
	Update(client Client) (Client, error)
	Delete(client Client) (Client, error)
	TotalCustomer(userID int) int
}

type repository struct {
	DB *gorm.DB
}

func NewClientRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(client Client) (Client, error) {
	err := r.DB.Create(&client).Error
	if err != nil {
		return client, err
	}

	return client, nil
}

func (r *repository) FindAll(s string, userID, page, perPage int) ([]Client, int, error) {
	var clients []Client
	var total int64

	sql := "SELECT * FROM clients WHERE deleted_at is null and user_id =  " + strconv.Itoa(userID)

	if s != "" {
		sql = fmt.Sprintf("%s AND fullname LIKE '%%%s%%'", sql, s)
	}

	r.DB.Raw(sql).Count(&total)

	sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, perPage, (page-1)*perPage)

	err := r.DB.Raw(sql).Scan(&clients).Error
	if err != nil {
		return clients, 0, err
	}

	return clients, int(total), nil
}

func (r *repository) FindById(id int) (Client, error) {
	var client Client
	err := r.DB.Where("id = ?", id).Find(&client).Error
	if err != nil {
		return client, err
	}

	return client, nil
}

func (r *repository) FindByEmail(email string, userID int) (Client, error) {
	var client Client
	err := r.DB.Where("email = ?", email).Where("user_id = ?", userID).Find(&client).Error
	if err != nil {
		return client, err
	}

	return client, nil
}

func (r *repository) Update(client Client) (Client, error) {
	err := r.DB.Save(&client).Error
	if err != nil {
		return client, err
	}

	return client, nil
}

func (r *repository) Delete(client Client) (Client, error) {
	err := r.DB.Delete(&client).Error
	if err != nil {
		return client, err
	}

	return client, nil
}

func (r *repository) TotalCustomer(userID int) int {
	var total int64

	sql := "SELECT * FROM clients WHERE user_id = " + strconv.Itoa(userID)
	r.DB.Raw(sql).Count(&total)

	return int(total)
}
