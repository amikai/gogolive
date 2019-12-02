package model

import (
	"errors"
	"reflect"
	"sync"
)

type User struct {
	BaseModel
	Name        string `db:"name"`
	Account     string `db:"account"`
	Password    string `db:"password"`
	ChannelName string `db:"channel_name"`
}

type IUserRepo interface {
	// Store User
	Store(User) error

	// Delete User
	Delete(User) error

	// Find the user by account. Accout name is unique, So just return one User.
	// If user not found, return nil, nil
	FindByAccount(account string) (*User, error)

	// Find the user by field. If user not found, retrun nil, nil.
	FilterByField(fieldName string, val interface{}) ([]User, error)

	// Update the user. if user is not existed, error will occur.
	Update(User) error

	// Get all user. if user not found, return nil, nil
	GetAll() ([]User, error)
}

// In memory user repo, which implement IUserRepo
type InMemoryUserRepo struct {
	AccountMap sync.Map
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		AccountMap: sync.Map{},
	}

}

func (repo *InMemoryUserRepo) FindByAccount(account string) (*User, error) {
	ret, ok := repo.AccountMap.Load(account)
	if !ok {
		return nil, nil
	}
	user := ret.(User)

	return &user, nil
}

func (repo *InMemoryUserRepo) Store(user User) error {
	_, accountExist := repo.AccountMap.Load(user.Account)
	if accountExist {
		return errors.New("The user account is duplicated")
	}

	repo.AccountMap.Store(user.Account, user)
	return nil
}

func (repo *InMemoryUserRepo) Delete(user User) error {
	repo.AccountMap.Delete(user.Account)
	return nil
}

func (repo *InMemoryUserRepo) Update(user User) error {
	_, accountExist := repo.AccountMap.Load(user.Account)
	if accountExist {
		repo.AccountMap.Store(user.Account, user)
		return nil
	}
	return errors.New("This user not exist")
}

func (repo *InMemoryUserRepo) FilterByField(fieldName string, fieldVal interface{}) ([]User, error) {

	v := reflect.ValueOf(User{})
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, errors.New("User type does not has this field")
	}
	if field.Type() != reflect.TypeOf(fieldVal) {
		return nil, errors.New("User field type and assign type not the same")
	}

	var users []User

	if fieldName == "Account" {
		user, ok := repo.AccountMap.Load(fieldVal)
		if ok {
			users = append(users, user.(User))
		}
		return users, nil
	}

	repo.AccountMap.Range(func(key, value interface{}) bool {
		user := value.(User)
		gotFieldVal := reflect.ValueOf(user).FieldByName(fieldName).Interface()
		if fieldVal == gotFieldVal {
			users = append(users, user)
		}
		return true
	})

	return users, nil
}

func (repo *InMemoryUserRepo) GetAll() ([]User, error) {
	var users []User
	repo.AccountMap.Range(func(key, value interface{}) bool {
		user := value.(User)
		users = append(users, user)
		return true
	})

	return users, nil
}
