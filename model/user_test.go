package model

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type InMemoryUserRepoTestSuite struct {
	suite.Suite
}

func (suite *InMemoryUserRepoTestSuite) TestStore() {
	user1 := User{Account: "test1"}
	user2 := User{Account: "test1"}
	suite.Run("store user", func() {
		userRepo := NewInMemoryUserRepo()
		userRepo.Store(user1)
		ret, ok := userRepo.AccountMap.Load(user1.Account)
		suite.True(ok)

		wantUser := user1
		gotUser, _ := ret.(User)
		suite.Equal(wantUser, gotUser)
	})

	suite.Run("user duplication", func() {
		var err error
		userRepo := NewInMemoryUserRepo()
		err = userRepo.Store(user1)
		suite.Nil(err)
		err = userRepo.Store(user2)
		suite.NotNil(err)
	})
}

func (suite *InMemoryUserRepoTestSuite) TestFindByAccount() {
	user := User{Account: "test"}
	suite.Run("user found", func() {
		userRepo := NewInMemoryUserRepo()
		userRepo.AccountMap.Store(user.Account, user)

		wantUser := user
		gotUser, err := userRepo.FindByAccount(user.Account)
		suite.Nil(err)
		suite.Equal(*gotUser, wantUser)
	})

	suite.Run("user not found", func() {
		userRepo := NewInMemoryUserRepo()
		user, _ := userRepo.FindByAccount(user.Account)
		suite.Nil(user)
	})
}

func (suite *InMemoryUserRepoTestSuite) TestDelete() {
	user := User{Account: "test1", Name: "name"}
	userRepo := NewInMemoryUserRepo()
	userRepo.AccountMap.Store(user.Account, user)

	userRepo.Delete(user)
	_, ok := userRepo.AccountMap.Load(user.Account)
	suite.False(ok)

}

func (suite *InMemoryUserRepoTestSuite) TestUpdate() {
	user := User{Account: "test1", Name: "name"}
	updatedUser := User{Account: "test1", Name: "updateName"}
	suite.Run("user not found", func() {
		userRepo := NewInMemoryUserRepo()
		err := userRepo.Update(user)
		suite.NotNil(err)
	})

	suite.Run("update user", func() {
		userRepo := NewInMemoryUserRepo()
		userRepo.AccountMap.Store(user.Account, user)

		err := userRepo.Update(updatedUser)
		suite.Nil(err)

		ret, _ := userRepo.AccountMap.Load(user.Account)
		gotUser, _ := ret.(User)
		wantUser := updatedUser
		suite.Equal(gotUser, wantUser)
	})
}

func (suite *InMemoryUserRepoTestSuite) TestGetAll() {
	user1 := User{Account: "test1"}
	user2 := User{Account: "test1"}
	expectedUserMap := map[string]User{
		user1.Account: user1,
		user2.Account: user2,
	}

	userRepo := NewInMemoryUserRepo()
	userRepo.AccountMap.Store(user1.Account, user1)
	userRepo.AccountMap.Store(user2.Account, user2)

	gotUsers, _ := userRepo.GetAll()
	gotUserMap := func() map[string]User {
		m := make(map[string]User)
		for _, user := range gotUsers {
			m[user.Account] = user
		}
		return m
	}()
	suite.Equal(expectedUserMap, gotUserMap)

}

func (suite *InMemoryUserRepoTestSuite) TestFilterByField() {
	user1 := User{Account: "test1", Name: "testName", ChannelName: "goodChannel"}
	user2 := User{Account: "test2", Name: "testName", ChannelName: "goodChannel"}
	user3 := User{Account: "test3", Name: "testName", ChannelName: "coolChannel"}
	initUserRepo := func() *InMemoryUserRepo {
		userRepo := NewInMemoryUserRepo()
		userRepo.AccountMap.Store(user1.Account, user1)
		userRepo.AccountMap.Store(user2.Account, user2)
		userRepo.AccountMap.Store(user3.Account, user3)
		return userRepo

	}

	suite.Run("User does not has this field", func() {
		userRepo := initUserRepo()
		_, err := userRepo.FilterByField("MAGIC", "123")
		suite.NotNil(err)

	})

	suite.Run("Wrong value type", func() {
		userRepo := initUserRepo()
		_, err := userRepo.FilterByField("Account", struct{}{})
		suite.NotNil(err)
	})

	suite.Run("Filter by Name", func() {
		userRepo := initUserRepo()
		gotUsers, _ := userRepo.FilterByField("Name", "testName")
		gotUserMap := func() map[string]User {
			m := make(map[string]User)
			for _, user := range gotUsers {
				m[user.Account] = user
			}
			return m
		}()

		expectedUserMap := map[string]User{
			user1.Account: user1,
			user2.Account: user2,
			user3.Account: user3,
		}

		suite.Equal(gotUserMap, expectedUserMap)

	})

	suite.Run("Filter by account", func() {

		suite.Run("Find user", func() {
			userRepo := initUserRepo()
			var gotUsers []User
			gotUsers, _ = userRepo.FilterByField("Account", "test1")
			suite.Equal(gotUsers[0], user1)
			gotUsers, _ = userRepo.FilterByField("Account", "test2")
			suite.Equal(gotUsers[0], user2)
			gotUsers, _ = userRepo.FilterByField("Account", "test3")
			suite.Equal(gotUsers[0], user3)

		})

		suite.Run("User not found", func() {
			userRepo := initUserRepo()
			users, _ := userRepo.FilterByField("Account", "not found")
			suite.Nil(users)
		})
	})

}

func TestInMemoryUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(InMemoryUserRepoTestSuite))
}
