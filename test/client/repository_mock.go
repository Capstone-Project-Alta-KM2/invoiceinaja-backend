package client

import (
	"errors"
	"invoiceinaja/domain/client"

	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	Mock mock.Mock
}

func NewRepositoryMock(mock mock.Mock) *RepositoryMock {
	return &RepositoryMock{mock}
}

func (r *RepositoryMock) Save(params client.Client) (client.Client, error) {
	argument := r.Mock.Called(params)
	if argument.Get(0) == nil {
		return params, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).(client.Client)
		return client, nil
	}
}

func (r *RepositoryMock) FindAll(s string, userID, page, perPage int) ([]client.Client, int, error) {
	argument := r.Mock.Called(s, userID, page, perPage)
	if argument.Get(0) == nil {
		return []client.Client{}, 0, errors.New("ada yang salah")
	} else {
		clients := argument.Get(0).([]client.Client)
		return clients, 0, nil
	}
}

func (r *RepositoryMock) FindById(id int) (client.Client, error) {
	argument := r.Mock.Called(id)
	if argument.Get(0) == nil {
		return client.Client{}, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).(client.Client)
		return client, nil
	}
}

func (r *RepositoryMock) FindByEmail(email string, userID int) (client.Client, error) {
	argument := r.Mock.Called(email, userID)
	if argument.Get(0) == nil {
		return client.Client{}, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).(client.Client)
		return client, nil
	}
}

func (r *RepositoryMock) Delete(params client.Client) (client.Client, error) {
	argument := r.Mock.Called(params)
	if argument.Get(0) == nil {
		return client.Client{}, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).(client.Client)
		return client, nil
	}
}

func (r *RepositoryMock) Update(params client.Client) (client.Client, error) {
	argument := r.Mock.Called(params)
	if argument.Get(0) == nil {
		return client.Client{}, errors.New("ada yang salah")
	} else {
		client := argument.Get(0).(client.Client)
		return client, nil
	}
}

func (r *RepositoryMock) TotalCustomer(userID int) int {
	argument := r.Mock.Called(userID)
	if argument.Get(0) == nil {
		return 0
	} else {
		return 12
	}
}
