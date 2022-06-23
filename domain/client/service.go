package client

type IService interface {
	AddClient(userID int, input InputAddClient) (Client, error)
	GetAll(clientID int, page int) ([]Client, int, int, error)
	GetByID(clientID int) (Client, error)
	UpdateClient(userID int, input InputUpdate) (Client, error)
	DeleteClient(clientID int) (Client, error)
}

type service struct {
	repository IRepository
}

func NewUserService(repository IRepository) *service {
	return &service{repository}
}

func (s *service) AddClient(userID int, input InputAddClient) (Client, error) {
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

func (s *service) GetAll(clientID int, page int) ([]Client, int, int, error) {
	perPage := 5
	clients, total, err := s.repository.FindAll(clientID, page, perPage)
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

func (s *service) UpdateClient(userID int, input InputUpdate) (Client, error) {
	client, errClient := s.repository.FindById(userID)
	if errClient != nil {
		return client, errClient
	}

	client.Fullname = input.Fullname
	client.Email = input.Email
	client.Address = input.Address
	client.City = input.City
	client.ZipCode = input.ZipCode
	client.Company = input.Company

	updatedUser, errUpdate := s.repository.Update(client)
	if errUpdate != nil {
		return updatedUser, errUpdate
	}

	return updatedUser, nil
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
