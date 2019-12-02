package service

import (
	"github.com/amikai/gogolive/model"
)

type Service struct {
	UserSerivce IUserService
	// TODO: if add new service, put here
}

func NewService(repo *model.Repo) *Service {
	return &Service{
		UserSerivce: NewUserService(repo),
	}
}
