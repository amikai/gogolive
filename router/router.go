package router

import (
	"github.com/amikai/gogolive/config"
	"github.com/amikai/gogolive/service"
	"github.com/gin-gonic/gin"
)

func Init(service *service.Service) *gin.Engine {
	if config.Conf.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.POST("/register", Register(service))
	r.POST("/signin", Signin(service))
	r.GET("/helloworld", HelloWorld())
	return r
}
