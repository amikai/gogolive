package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/amikai/gogolive/config"
	"github.com/amikai/gogolive/model"
	"github.com/amikai/gogolive/router"
	"github.com/amikai/gogolive/service"
)

func initLog() {
	if l, err := log.ParseLevel(config.Conf.LogLevel); err == nil {
		log.SetLevel(l)
		log.SetReportCaller(l == log.DebugLevel)
	}
}

func startWebService() {
	service := service.NewService(model.NewInMemoryRepo())
	r := router.Init(service)

	port := config.Conf.WebServer.Port
	log.Infof("Web service listen on %s", port)
	r.Run(port)
	log.Fatal(http.ListenAndServe(port, r))
}

func main() {
	// Use default config
	err := config.LoadConf("")
	if err != nil {
		log.Printf("Load yaml config file error: '%v'", err)
		return
	}

	initLog()
	startWebService()
}
