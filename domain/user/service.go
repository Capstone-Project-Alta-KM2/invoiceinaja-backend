package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type IService interface {
	Register(input InputRegister) (User, error)
	Login(input InputLogin) (User, error)
	GetUserById(id int) (User, error)
	IsEmailAvailable(input InputCheckEmail) (bool, error)
}

type service struct {
	repository IRepository
}

func NewUserService(repository IRepository) *service {
	return &service{repository}
}

func (s *service) Register(input InputRegister) (User, error) {
	var newUser User

	//enkripsi password
	passwordHash, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if errHash != nil {
		return newUser, errHash
	}

	//tangkap nilai dari inputan
	newUser.Fullname = input.Fullname
	newUser.Email = input.Email
	newUser.BusinessName = input.BusinessName
	newUser.Password = string(passwordHash)
	newUser.Role = "user"
	newUser.Avatar = "images/default_user.png"

	//save data yang sudah dimapping kedalam struct Mahasiswa
	user, err := s.repository.Save(newUser)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) IsEmailAvailable(input InputCheckEmail) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) Login(input InputLogin) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	//cek jika user tidak ada
	if user.ID == 0 {
		return user, errors.New("no user found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) GetUserById(id int) (User, error) {
	user, err := s.repository.FindById(id)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found")
	}

	return user, nil
}
