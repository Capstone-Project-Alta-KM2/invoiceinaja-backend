package user

import (
	"testing"
	"time"

	"invoiceinaja/domain/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

var repoUser = &RepositoryMock{Mock: mock.Mock{}}
var serviceUser = user.NewUserService(repoUser)

func TestRegister(t *testing.T) {
	var input = user.InputRegister{
		Fullname:     "fahmi",
		Email:        "fahmi@gmail.com",
		NoTlpn:       "08123",
		BusinessName: "asda",
		Password:     "password",
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	usr := user.User{
		ID:           0,
		Fullname:     "fahmi",
		Email:        "fahmi@gmail.com",
		NoTlpn:       "08123",
		BusinessName: "asda",
		Password:     string(passwordHash),
		Role:         "user",
		Avatar:       "images/default_user.png",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}

	repoUser.Mock.On("FindByEmail", input.Email).Return(usr)
	repoUser.Mock.On("Save", mock.AnythingOfType("user.User")).Return(usr)

	result, err := serviceUser.Register(input)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, usr.ID, result.ID)
	assert.True(t, bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(input.Password)) == nil)
}
