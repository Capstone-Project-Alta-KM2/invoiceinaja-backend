package user

import (
	"errors"

	"invoiceinaja/domain/user"

	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	Mock mock.Mock
}

func NewRepositoryMock(mock mock.Mock) *RepositoryMock {
	return &RepositoryMock{mock}
}

func (r *RepositoryMock) Save(user user.User) (user.User, error) {

	argument := r.Mock.Called(user)
	if argument.Get(0) == nil {
		return user, errors.New("ada yang salah")
	} else {
		// newUser := argument.Get(0)
		return user, nil
	}
}

func (r *RepositoryMock) FindById(id int) (user.User, error) {
	argument := r.Mock.Called(id)
	if argument.Get(0) == nil {
		return user.User{}, errors.New("ada yang salah")
	} else {
		user := argument.Get(0).(user.User)
		return user, nil
	}
}

func (r *RepositoryMock) FindByEmail(email string) (user.User, error) {
	var user user.User
	argument := r.Mock.Called(email)
	if argument.Get(0) == nil {
		return user, errors.New("ada yang salah")
	} else {
		// user := argument.Get(0).(user.User)
		return user, nil
	}
}

func (r *RepositoryMock) Update(user user.User) (user.User, error) {
	argument := r.Mock.Called(user)
	if argument.Get(0) == nil {
		return user, errors.New("ada yang salah")
	} else {
		// newUser := argument.Get(0).(user.User)
		return user, nil
	}
}
func (r *RepositoryMock) UpdateByID(user user.User) (user.User, error) {
	argument := r.Mock.Called(user)
	if argument.Get(0) == nil {
		return user, errors.New("ada yang salah")
	} else {
		// newUser := argument.Get(0).(user.User)
		return user, nil
	}
}

func (r *RepositoryMock) DeleteUser(user user.User) (user.User, error) {
	argument := r.Mock.Called(user)
	if argument.Get(0) == nil {
		return user, errors.New("ada yang salah")
	} else {
		return user, nil
	}

}
