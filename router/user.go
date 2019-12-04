package router

import (
	"encoding/json"
	"net/http"

	"github.com/amikai/gogolive/model"
	"github.com/amikai/gogolive/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/qri-io/jsonschema"
)

func jsonValidate(schema []byte) bool {
	rs := &jsonschema.RootSchema{}
	if err := json.Unmarshal(schema, rs); err != nil {
		return false
	}
	return true
}

func reportStatus(c *gin.Context, code int, status, mes string) {
	c.JSON(code, gin.H{"status": status, "mes": mes})
	return
}

func Signin(s *service.Service) gin.HandlerFunc {
	signinForm := struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	return func(c *gin.Context) {
		var err error
		err = c.ShouldBindJSON(&signinForm)
		if err != nil {
			reportStatus(c, http.StatusBadRequest, "error", "wrong json format")
			return
		}

		err = s.UserSerivce.VerifyPassword(model.User{
			Account:  signinForm.Account,
			Password: signinForm.Password,
		})

		if err != nil {
			reportStatus(c, http.StatusBadRequest, "error", err.Error())
			return
		}
		reportStatus(c, http.StatusOK, "ok", "successfully signin")
	}
}

func Register(s *service.Service) gin.HandlerFunc {
	registerForm := struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	return func(c *gin.Context) {
		var err error
		// bind json form
		err = c.ShouldBindBodyWith(&registerForm, binding.JSON)
		if err != nil {
			reportStatus(c, http.StatusBadRequest, "error", "wrong json format")
			return
		}

		//  register user
		err = s.UserSerivce.Register(model.User{
			Account:  registerForm.Account,
			Password: registerForm.Password,
		})

		if err != nil {
			reportStatus(c, http.StatusInternalServerError, "error", err.Error())
			return
		}

		reportStatus(c, http.StatusOK, "ok", "successfully register")
	}
}

func HelloWorld() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hello": "world",
		})
	}
}
