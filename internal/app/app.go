package app

import (
	"log"

	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
)

func Run() {
	db, err := repository.InitDB(repository.Config{
		Host:     "localhost",
		Port:     "8080",
		Username: "sqlite",
		Password: "qwerty",
		DBName:   "forumDB",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := delivery.NewHandler(services)

	server := new(server.Server)
	if err := server.ServerRun("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error ocured while running http server: %s", err.Error())
	}
}
