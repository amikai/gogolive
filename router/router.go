package router

import (
	"github.com/amikai/gogolive/service"
	"github.com/gin-gonic/gin"
)

func Init(service *service.Service) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/register", Register(service))
	r.POST("/signin", Signin(service))
	r.GET("/helloworld", HelloWorld())
	return r
}
