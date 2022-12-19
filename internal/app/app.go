package app

import (
	"log"

	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
)

func Run() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := delivery.NewHandler(services)

	server := new(server.Server)
	if err := server.ServerRun("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error ocured while running http server: %s", err.Error())
	}
}
