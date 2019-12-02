package main

import (
	"log"
	"net/http"

	"github.com/amikai/gogolive/model"
	"github.com/amikai/gogolive/router"
	"github.com/amikai/gogolive/service"
)

func main() {
	service := service.NewService(model.NewInMemoryRepo())
	r := router.Init(service)
	log.Fatal(http.ListenAndServe(":8080", r))
}
