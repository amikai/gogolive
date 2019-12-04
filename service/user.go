package service

import (
	"errors"

	"github.com/amikai/gogolive/model"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Register(user model.User) error
	VerifyPassword(user model.User) (bool, error)
}

type UserService struct {
	Repo *model.Repo
}

func NewUserService(repo *model.Repo) *UserService {
	return &UserService{
		Repo: repo,
	}
}

func hashAndSalt(password []byte) (string, bool) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return "", false
	}
	return string(hash), true
}

func comparePassword(hashedPassword string, plainPassword []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		return false
	}
	return true
}

// VerifyPassword verify password when login
func (service *UserService) VerifyPassword(user model.User) (bool, error) {
	var err error
	plainPassword := []byte(user.Password)
	userInStore, err := service.Repo.UserRepo.FindByAccount(user.Account)
	if userInStore == nil {
		return false, err
	}
	hashedPassword := []byte(userInStore.Password)
	err = bcrypt.CompareHashAndPassword(hashedPassword, plainPassword)
	if err != nil {
		return false, err
	}
	return true, nil
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
