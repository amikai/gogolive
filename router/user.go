package router

import (
	"encoding/json"
	"net/http"

	"github.com/amikai/gogolive/config"
	"github.com/amikai/gogolive/model"
	"github.com/amikai/gogolive/service"
	"github.com/dgrijalva/jwt-go"
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

type JWTClaims struct {
	Account string `json:"account"`
	Age     int    `json:"age"`
	jwt.StandardClaims
}

func Signin(s *service.Service) gin.HandlerFunc {
	signinForm := struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	jwtKey := config.Conf.JWT.Key
	jwtAge := config.Conf.JWT.Age

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

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaims{
			signinForm.Account,
			jwtAge,
			jwt.StandardClaims{},
		})
		tokenString, err := token.SignedString([]byte(jwtKey))
		if err != nil {
			reportStatus(c, http.StatusInternalServerError, "error", err.Error())
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			MaxAge:   jwtAge,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

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
