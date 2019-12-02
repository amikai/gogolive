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
			c.JSON(http.StatusBadRequest,
				gin.H{
					"status": "error",
					"mes":    "wrong json format",
				})
			return
		}

		//  register user
		err = s.UserSerivce.Register(model.User{
			Account:  registerForm.Account,
			Password: registerForm.Password,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"status": "error",
					"mes":    err.Error(),
				})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func HelloWorld() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hello": "world",
		})
	}
}
