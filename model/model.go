package model

import (
	"time"
)

type BaseModel struct {
	Id        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

type Repo struct {
	UserRepo IUserRepo
	// TODO: If add new repo, put here
}

func NewInMemoryRepo() *Repo {
	return &Repo{
		UserRepo: NewInMemoryUserRepo(),
	}
}
