package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/amikai/gogolive/config"
	"github.com/amikai/gogolive/service"
	"github.com/amikai/gogolive/service/mock"
	"github.com/dgrijalva/jwt-go"
	"github.com/gavv/httpexpect"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	// setup gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/helloWorld", HelloWorld())

	// setup test server
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	schema := `
	{
	  "type": "object",
	  "properties": {
		"hello": {
		  "type": "string"
		}
	  },
	  "required": [
		"hello"
	  ]
	}`
	repos := e.GET("/helloWorld").Expect().
		Status(http.StatusOK).
		JSON()

	repos.Schema(schema)
	repos.Path(`$["hello"]`).Equal("world")
}

func TestRegister(t *testing.T) {
	// setup go mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	requestPath := "/register"

	setUpMockAndRouter := func() (*mock.MockIUserService, *httptest.Server) {
		mockUserService := mock.NewMockIUserService(mockCtrl)

		gin.SetMode(gin.DebugMode)
		router := gin.New()
		router.POST(requestPath, Register(&service.Service{UserSerivce: mockUserService}))
		server := httptest.NewServer(router)
		return mockUserService, server
	}

	newRequest := func(t *testing.T, serverURL string, path string) *httpexpect.Request {
		return httpexpect.New(t, serverURL).POST(path)
	}

	t.Run("missing field", func(t *testing.T) {

		mockUserService, server := setUpMockAndRouter()
		defer server.Close()

		mockUserService.EXPECT().Register(gomock.Any()).Return(nil).Times(0)
		newRequest(t, server.URL, requestPath).
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("correct form", func(t *testing.T) {
		mockUserService, server := setUpMockAndRouter()
		defer server.Close()

		mockUserService.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		respJSON := newRequest(t, server.URL, requestPath).
			WithJSON(map[string]interface{}{
				"account":  "account",
				"password": "pass",
			}).
			Expect().Status(http.StatusOK).JSON()
		respJSON.Path("$.status").Equal("ok")
	})

	t.Run("additional field", func(t *testing.T) {
		mockUserService, server := setUpMockAndRouter()
		defer server.Close()

		mockUserService.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		respJSON := newRequest(t, server.URL, requestPath).
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]string{
				"account":          "account",
				"password":         "pass",
				"additional field": "test",
			}).
			Expect().Status(http.StatusOK).JSON()
		respJSON.Path("$.status").Equal("ok")
	})

	t.Run("service register faile", func(t *testing.T) {
		mockUserService, server := setUpMockAndRouter()
		defer server.Close()

		mockUserService.EXPECT().Register(gomock.Any()).Return(errors.New("mock")).Times(1)
		respJSON := newRequest(t, server.URL, requestPath).
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]string{
				"account":  "account",
				"password": "pass",
			}).
			Expect().Status(http.StatusInternalServerError).JSON()
		respJSON.Path("$.status").Equal("error")
	})

}

func TestSignin(t *testing.T) {
	// setup go mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	requestPath := "/signin"

	setUpMockAndRouter := func() (*mock.MockIUserService, *httptest.Server) {
		mockUserService := mock.NewMockIUserService(mockCtrl)

		gin.SetMode(gin.DebugMode)
		router := gin.New()
		router.POST(requestPath, Signin(&service.Service{UserSerivce: mockUserService}))
		server := httptest.NewServer(router)
		return mockUserService, server
	}

	newRequest := func(t *testing.T, serverURL string, path string) *httpexpect.Request {
		return httpexpect.New(t, serverURL).POST(path)
	}

	t.Run("invalid field", func(t *testing.T) {
		mockUserService, server := setUpMockAndRouter()
		defer server.Close()

		mockUserService.EXPECT().VerifyPassword(gomock.Any()).Return(nil).Times(0)
		respJSON := newRequest(t, server.URL, requestPath).
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{}).
			Expect().
			Status(http.StatusBadRequest).JSON()
		respJSON.Path("$.status").Equal("error")
	})

	t.Run("successful signin", func(t *testing.T) {
		mockUserService, server := setUpMockAndRouter()
		defer server.Close()

		mockUserService.EXPECT().VerifyPassword(gomock.Any()).Return(nil).Times(1)
		respJSON := newRequest(t, server.URL, requestPath).
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"account":  "account",
				"password": "password",
			}).
			Expect().
			Status(http.StatusOK).JSON()
		respJSON.Path("$.status").Equal("ok")
	})

	t.Run("verify password failed", func(t *testing.T) {
		mockUserService, server := setUpMockAndRouter()
		defer server.Close()

		mockUserService.EXPECT().VerifyPassword(gomock.Any()).Return(errors.New("mock")).Times(1)
		respJSON := newRequest(t, server.URL, requestPath).
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"account":  "account",
				"password": "password",
			}).
			Expect().
			Status(http.StatusBadRequest).JSON()
		respJSON.Path("$.status").Equal("error")
	})

	t.Run("test token", func(t *testing.T) {
		config.Conf.JWT.Key = "JWT_KEY"
		config.Conf.JWT.Age = 6666 // magic number
		mockUserService, server := setUpMockAndRouter()
		defer server.Close()

		mockUserService.EXPECT().VerifyPassword(gomock.Any()).Return(nil).Times(1)
		resp := newRequest(t, server.URL, requestPath).
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"account":  "account",
				"password": "password",
			}).
			Expect().
			Status(http.StatusOK)

		gotMaxAge := config.Conf.JWT.Age

		cookie := resp.Cookie("token")
		tokenString := cookie.Value().Raw()
		jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			claims := token.Claims.(*JWTClaims)
			assert.Equal(t, claims.Account, "account")
			assert.Equal(t, claims.Age, gotMaxAge)
			return nil, nil
		})
		cookie.MaxAge().Equal(time.Duration(gotMaxAge) * time.Second)

		respJSON := resp.JSON()
		respJSON.Path("$.status").Equal("ok")

	})
}
