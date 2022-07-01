package client

import (
	"errors"
	"invoiceinaja/utils"
)

type IService interface {
	AddClient(userID int, input InputAddClient) (Client, error)
	IsEmailClientAvailable(input string, userID int) (bool, error)
	GetAll(search string, clientID, page int) ([]Client, int, int, error)
	GetByID(clientID int) (Client, error)
	UpdateClient(input InputUpdate, userID, clientID int) (Client, error)
	DeleteClient(clientID int) (Client, error)
}

type service struct {
	repository IRepository
}

func NewUserService(repository IRepository) *service {
	return &service{repository}
}

func (s *service) IsEmailClientAvailable(input string, userID int) (bool, error) {
	client, err := s.repository.FindByEmail(input, userID)
	if err != nil {
		return false, err
	}

	if client.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) AddClient(userID int, input InputAddClient) (Client, error) {

	clientExist, errExist := s.repository.FindByEmail(input.Email, userID)
	if clientExist.ID != 0 {
		return clientExist, errExist
	}

	var data Client

	data.Fullname = input.Fullname
	data.Email = input.Email
	data.Address = input.Address
	data.City = input.City
	data.ZipCode = input.ZipCode
	data.Company = input.Company
	data.UserID = userID

	client, err := s.repository.Save(data)
	if err != nil {
		return client, err
	}

	return client, nil
}

func (s *service) GetAll(search string, clientID, page int) ([]Client, int, int, error) {
	perPage := 5
	clients, total, err := s.repository.FindAll(search, clientID, page, perPage)
	if err != nil {
		return clients, 0, 0, err
	}

	return clients, total, perPage, nil
}

func (s *service) GetByID(clientID int) (Client, error) {
	client, err := s.repository.FindById(clientID)
	if err != nil {
		return client, err
	}

	return client, nil
}

func (s *service) UpdateClient(input InputUpdate, userID, clientID int) (Client, error) {

	client, errItem := s.repository.FindById(clientID)
	if errItem != nil {
		return client, errItem
	}
	if client.UserID != userID || client.UserID == 0 {
		return Client{}, errors.New("bad request")
	}

	client.ID = clientID
	client.Fullname = input.Fullname
	client.Email = input.Email
	client.Address = input.Address
	client.City = input.City
	client.ZipCode = input.ZipCode
	client.Company = input.Company
	client.UserID = userID

	updatedItem, errUpdate := s.repository.Update(client)
	if errUpdate != nil {
		return updatedItem, errUpdate
	}

	return updatedItem, nil
}

func (s *service) DeleteClient(clientID int) (Client, error) {
	client, err := s.repository.FindById(clientID)
	if err != nil {
		return client, err
	}

	deleteClient, errDel := s.repository.Delete(client)
	if errDel != nil {
		return deleteClient, errDel
	}

	return deleteClient, nil
}

func Mapping(lines [][]string) []InputAddClient {
	var clients []InputAddClient
	var email []string // nampung email

	// Loop through lines & turn into object
	for i, line := range lines {
		if i == 0 {
			continue
		}

		// mapping data
		data := InputAddClient{
			Fullname: line[0],
			Email:    line[1],
			Address:  line[2],
			City:     line[3],
			ZipCode:  line[4],
			Company:  line[5],
		}

		// if email sudah ada pada slice email
		// email tidak akan diappend kedalam slice email
		isExist := utils.Contains(email, data.Email)
		switch isExist {
		case true:
			continue
		case false:
			email = append(email, data.Email)
			clients = append(clients, InputAddClient{
				Fullname: data.Fullname,
				Email:    data.Email,
				Address:  data.Address,
				City:     data.City,
				ZipCode:  data.ZipCode,
				Company:  data.Company,
			})
		}

	}

	return clients
}
