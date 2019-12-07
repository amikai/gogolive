package service

import (
	"errors"

	"github.com/amikai/gogolive/model"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Register(user model.User) error
	VerifyPassword(user model.User) error
}

type UserService struct {
	Repo *model.Repo
}

func NewUserService(repo *model.Repo) *UserService {
	return &UserService{
		Repo: repo,
	}
}

// VerifyPassword verify password when login
func (service *UserService) VerifyPassword(user model.User) error {
	var err error
	plainPassword := []byte(user.Password)
	userInStore, err := service.Repo.UserRepo.FindByAccount(user.Account)
	if userInStore == nil {
		return errors.New("user not found")
	}
	hashedPassword := []byte(userInStore.Password)
	err = bcrypt.CompareHashAndPassword(hashedPassword, plainPassword)
	return err
}

// register account
func (service *UserService) Register(user model.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Password encryption failed")
	}
	user.Password = string(password)
	err = service.Repo.UserRepo.Store(user)

	if err != nil {
		return errors.New("User register failed")
	}

	return nil
}

// or somthing high level operation
