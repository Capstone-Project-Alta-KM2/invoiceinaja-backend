package database

import "invoiceinaja/domain/user"

type Model struct {
	Model interface{}
}

func RegisterModel() []Model {
	return []Model{
		{Model: user.User{}},
	}
}
