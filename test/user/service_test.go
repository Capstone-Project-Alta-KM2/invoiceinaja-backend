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

func TestLogin(t *testing.T) {
	var input = user.InputLogin{
		Email:    "hadifah3@gmail.com",
		Password: "password",
	}
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	usr := user.User{
		ID:           1,
		Fullname:     "fahmi",
		Email:        input.Email,
		NoTlpn:       "08123",
		BusinessName: "asda",
		Password:     string(passwordHash),
		Role:         "user",
		Avatar:       "images/default_user.png",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}

	repoUser.Mock.On("FindByEmail", input.Email).Return(usr)

	result, err := serviceUser.Login(input)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, usr.Email, result.Email)
}

func TestGetUserById(t *testing.T) {
	user := user.User{
		ID:           12,
		Fullname:     "fahmi",
		Email:        "fahmi@gmail.com",
		NoTlpn:       "0852136414",
		Password:     "password",
		BusinessName: "toko online",
		Role:         "user",
		Avatar:       "images/default_user.png",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	repoUser.Mock.On("FindById", 12).Return(user)

	result, err := serviceUser.GetUserById(12)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Fullname, result.Fullname)
}

func TestIsEmailAvailable(t *testing.T) {
	var input = user.InputCheckEmail{Email: "fahmi@gmail.com"}
	user := user.User{
		ID:           12,
		Fullname:     "fahmi",
		Email:        "fahmi@gmail.com",
		NoTlpn:       "0852136414",
		Password:     "password",
		BusinessName: "toko online",
		Role:         "user",
		Avatar:       "images/default_user.png",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	repoUser.Mock.On("FindByEmail", input.Email).Return(user)

	result, err := serviceUser.IsEmailAvailable(input)
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

func TestSaveAvatar(t *testing.T) {
	path := "images/avatar.png"
	user := user.User{
		ID:           12,
		Fullname:     "fahmi",
		Email:        "fahmi@gmail.com",
		NoTlpn:       "0852136414",
		Password:     "password",
		BusinessName: "toko online",
		Role:         "user",
		Avatar:       "images/default_user.png",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	repoUser.Mock.On("FindById", user.ID).Return(user)
	repoUser.Mock.On("Update", mock.AnythingOfType("user.User")).Return(user)

	res, err := serviceUser.SaveAvatar(user.ID, path)
	user = res

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, user.Avatar, res.Avatar)
}

func TestUpdateUser(t *testing.T) {
	input := user.InputUpdate{
		Fullname:     "fahmi hadi",
		Email:        "hadifah3@gmail.com",
		NoTlpn:       "085214796",
		BusinessName: "ojol",
	}
	user := user.User{
		ID:           12,
		Fullname:     "fahmi hadi",
		Email:        "hadifah3@gmail.com",
		NoTlpn:       "085214796",
		Password:     "password",
		BusinessName: "ojol",
		Role:         "user",
		Avatar:       "images/default_user.png",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	repoUser.Mock.On("FindById", user.ID).Return(user)
	repoUser.Mock.On("Update", user).Return(user)

	res, err := serviceUser.UpdateUser(user.ID, input)
	user = res

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, user.ID, res.ID)
	assert.Equal(t, user.Email, res.Email)
}

func TestChangePassword(t *testing.T) {

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("old_password"), bcrypt.MinCost)
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

	repoUser.Mock.On("FindById", usr.ID).Return(usr)
	repoUser.Mock.On("Update", mock.AnythingOfType("user.User")).Return(usr)

	res, err := serviceUser.ChangePassword(user.InputChangePassword{
		OldPassword: "old_password",
		NewPassword: "new_password",
	}, usr.ID)

	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestResetPassword(t *testing.T) {
	var input = user.InputCheckEmail{Email: "fahmi@gmail.com"}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
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
	repoUser.Mock.On("Update", mock.AnythingOfType("user.User")).Return(usr)

	res, err := serviceUser.ResetPassword(input)

	assert.Nil(t, err)
	assert.NotNil(t, res)
}
