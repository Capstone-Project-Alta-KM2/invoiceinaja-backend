package client

import (
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/user"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var repoClient = &RepositoryMock{Mock: mock.Mock{}}
var serviceClient = client.NewUserService(repoClient)

func TestAddClient(t *testing.T) {

	inp := client.InputAddClient{
		Fullname: "Aldi",
		Email:    "aldi@gmail.com",
		Address:  "bandung",
		City:     "bandung",
		ZipCode:  "41523",
		Company:  "toko sembako",
	}

	clientData := client.Client{
		ID:        1,
		Fullname:  inp.Fullname,
		Email:     inp.Email,
		Address:   inp.Address,
		City:      inp.City,
		ZipCode:   inp.ZipCode,
		Company:   inp.Company,
		UserID:    1,
		User:      user.User{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}

	repoClient.Mock.On("FindByEmail", inp.Email, 1).Return(clientData)
	repoClient.Mock.On("Save", mock.AnythingOfType("client.Client")).Return(clientData)

	res, err := serviceClient.AddClient(1, inp)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, clientData.ID, res.ID)
}

func TestIsEmailClientAvailable(t *testing.T) {
	inp := client.InputAddClient{
		Fullname: "Aldi",
		Email:    "aldi@gmail.com",
		Address:  "bandung",
		City:     "bandung",
		ZipCode:  "41523",
		Company:  "toko sembako",
	}

	clientData := client.Client{
		ID:        1,
		Fullname:  inp.Fullname,
		Email:     inp.Email,
		Address:   inp.Address,
		City:      inp.City,
		ZipCode:   inp.ZipCode,
		Company:   inp.Company,
		UserID:    1,
		User:      user.User{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}

	repoClient.Mock.On("FindByEmail", inp.Email, 1).Return(clientData)

	res, err := serviceClient.IsEmailClientAvailable("aldi@gmail.com", clientData.UserID)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.True(t, !res)
}

func TestGetAll(t *testing.T) {
	clientsData := []client.Client{
		{
			ID:        1,
			Fullname:  "client1",
			Email:     "client1@gmail.com",
			Address:   "bandung",
			City:      "bandung",
			ZipCode:   "41523",
			Company:   "toko sembako",
			UserID:    1,
			User:      user.User{},
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		{
			ID:        2,
			Fullname:  "client2",
			Email:     "client2@gmail.com",
			Address:   "bandung",
			City:      "bandung",
			ZipCode:   "41523",
			Company:   "toko sembako",
			UserID:    1,
			User:      user.User{},
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
	}

	repoClient.Mock.On("FindAll", "", 1, 0, 5).Return(clientsData)
	clients, _, _, err := serviceClient.GetAll("", 1, 0)

	assert.Nil(t, err)
	assert.NotNil(t, clients)
}

func TestGetByEmail(t *testing.T) {
	clientData := client.Client{
		ID:        1,
		Fullname:  "client2",
		Email:     "client2@gmail.com",
		Address:   "bandung",
		City:      "bandung",
		ZipCode:   "41523",
		Company:   "toko sembako",
		UserID:    1,
		User:      user.User{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}

	repoClient.Mock.On("FindByEmail", "client2@gmail.com", 1).Return(clientData)
	res, err := serviceClient.GetByEmail("client2@gmail.com", 1)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, clientData.ID, res.ID)
}

func TestGetByID(t *testing.T) {
	clientData := client.Client{
		ID:        1,
		Fullname:  "client2",
		Email:     "client2@gmail.com",
		Address:   "bandung",
		City:      "bandung",
		ZipCode:   "41523",
		Company:   "toko sembako",
		UserID:    1,
		User:      user.User{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}

	repoClient.Mock.On("FindById", clientData.ID).Return(clientData)
	res, err := serviceClient.GetByID(1)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, clientData.ID, res.ID)
}

func TestUpdateClient(t *testing.T) {
	inp := client.InputUpdate{
		Fullname: "client2",
		Email:    "client2@gmail.com",
		Address:  "bandung",
		City:     "bandung",
		ZipCode:  "41523",
		Company:  "toko sembako",
		User: user.User{
			ID: 1,
		},
	}

	clientData := client.Client{
		ID:        1,
		Fullname:  "client2",
		Email:     "client2@gmail.com",
		Address:   "bandung",
		City:      "bandung",
		ZipCode:   "41523",
		Company:   "toko sembako",
		UserID:    1,
		User:      user.User{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}

	repoClient.Mock.On("FindById", clientData.ID).Return(clientData)
	repoClient.Mock.On("Update", mock.AnythingOfType("client.Client")).Return(clientData)

	res, err := serviceClient.UpdateClient(inp, 1, 1)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, clientData.ID, res.ID)
}

func TestDeleteClient(t *testing.T) {
	clientData := client.Client{
		ID:        1,
		Fullname:  "client2",
		Email:     "client2@gmail.com",
		Address:   "bandung",
		City:      "bandung",
		ZipCode:   "41523",
		Company:   "toko sembako",
		UserID:    1,
		User:      user.User{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}

	repoClient.Mock.On("FindById", clientData.ID).Return(clientData)
	repoClient.Mock.On("Delete", mock.AnythingOfType("client.Client")).Return(clientData)

	res, err := serviceClient.DeleteClient(clientData.ID)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, clientData.ID, res.ID)
}

func TestTotalCustomer(t *testing.T) {
	repoClient.Mock.On("TotalCustomer", 1).Return(12)
	res := serviceClient.TotalCustomer(1)

	assert.NotNil(t, res)

}
