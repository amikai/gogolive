package service

import (
	"errors"
	"testing"

	"github.com/amikai/gogolive/model"
	"github.com/amikai/gogolive/model/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestVerifyPassword(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stringType := gomock.AssignableToTypeOf("")
	t.Run("user found", func(t *testing.T) {
		mockUserRepo := mock.NewMockIUserRepo(mockCtrl)
		repo := model.Repo{UserRepo: mockUserRepo}
		userService := NewUserService(&repo)

		// if you want to change hash function, change the hashedPassword
		plainPassword := "plainPassword"
		hashedPassword := "$2a$10$WU.qv8CLG3G.VnvQ/M3O2OrBaNf0VrZtC1YhEdckqWn4YtvQxs9j2"
		mockUserRepo.EXPECT().FindByAccount(stringType).Return(&model.User{Password: hashedPassword}, nil).Times(1)
		err := userService.VerifyPassword(model.User{Password: plainPassword})
		assert.NoError(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo := mock.NewMockIUserRepo(mockCtrl)
		repo := model.Repo{UserRepo: mockUserRepo}
		userService := NewUserService(&repo)

		mockUserRepo.EXPECT().FindByAccount(stringType).Return(nil, nil).Times(1)
		err := userService.VerifyPassword(model.User{})
		assert.Error(t, err)
	})

}

func TestRegister(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userType := gomock.AssignableToTypeOf(model.User{})

	t.Run("user register error", func(t *testing.T) {
		mockUserRepo := mock.NewMockIUserRepo(mockCtrl)
		repo := model.Repo{UserRepo: mockUserRepo}
		userService := NewUserService(&repo)
		mockUserRepo.EXPECT().Store(userType).Return(nil).Times(1)
		err := userService.Register(model.User{})
		assert.NoError(t, err)
	})

	t.Run("user register success", func(t *testing.T) {
		mockUserRepo := mock.NewMockIUserRepo(mockCtrl)
		repo := model.Repo{UserRepo: mockUserRepo}
		userService := NewUserService(&repo)
		mockUserRepo.EXPECT().Store(userType).Return(errors.New("")).Times(1)
		err := userService.Register(model.User{})
		assert.Error(t, err)
	})
}
