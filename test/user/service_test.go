package user

import (
	"fmt"
	"testing"
	"time"

	"invoiceinaja/domain/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	// passwordHash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	user := user.User{
		ID:           0,
		Fullname:     "fahmi",
		Email:        "fahmi@gmail.com",
		NoTlpn:       "08123",
		BusinessName: "asda",
		Password:     "password",
		Role:         "user",
		Avatar:       "images/default_user.png",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}

	repoUser.Mock.On("FindByEmail", input.Email).Return(user)
	repoUser.Mock.On("Save", user).Return(user)

	result, err := serviceUser.Register(input)

	fmt.Println(result)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Fullname, result.Fullname)
}
