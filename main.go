package main

import (
	"net"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/amikai/gogolive/config"
	"github.com/amikai/gogolive/model"
	"github.com/amikai/gogolive/protocol/rtmp"
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
	log.Fatal(http.ListenAndServe(port, r))
}

func startRTMPServer() {
	port := config.Conf.RtmpServer.Port
	rtmpListen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	rtmpServer := rtmp.NewRtmpServer()
	log.Infof("Rtmp listen on %s", port)
	rtmpServer.Serve(rtmpListen)
}

func main() {
	// Use default config
	err := config.LoadConf("")
	if err != nil {
		log.Printf("Load yaml config file error: '%v'", err)
		return
	}

	initLog()

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		startWebService()
		wg.Done()
	}()
	go func() {
		startRTMPServer()
		wg.Done()
	}()

	wg.Wait()
}
