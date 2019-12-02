package service

import (
	"errors"

	"github.com/amikai/gogolive/model"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Register(user model.User) error
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
