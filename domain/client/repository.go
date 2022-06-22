package client

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type IRepository interface {
	Save(client Client) (Client, error)
	FindAll(userID int, page, perPage int) ([]Client, int, error)
	FindById(id int) (Client, error)
	FindByEmail(email string) (Client, error)
	Update(client Client) (Client, error)
	Delete(client Client) (Client, error)
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

func (r *repository) FindAll(userID int, page, perPage int) ([]Client, int, error) {
	//page, _ := strconv.Atoi(c.Query("page", "1"))
	//perPage := 9
	//var total int64
	var clients []Client
	var total int64

	sql := "SELECT * FROM clients WHERE user_id = " + strconv.Itoa(userID)

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

func (r *repository) FindByEmail(email string) (Client, error) {
	var client Client

	err := r.DB.Where("email = ?", email).Find(&client).Error
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
