package client

import "invoiceinaja/domain/user"

type InputAddClient struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Address  string `json:"address" binding:"required"`
	City     string `json:"city" binding:"required"`
	ZipCode  string `json:"zip_code" binding:"required"`
	Company  string `json:"company" binding:"required"`
}

type InputUpdate struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Address  string `json:"address" binding:"required"`
	City     string `json:"city" binding:"required"`
	ZipCode  string `json:"zip_code" binding:"required"`
	Company  string `json:"company" binding:"required"`
	User     user.User
}
