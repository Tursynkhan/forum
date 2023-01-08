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
		Username: "sqlite3",
		DBName:   "forumDB.db",
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	if err := repository.CreateTables(db); err != nil {
		log.Println(err)
	}
	defer db.Close()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := delivery.NewHandler(services)

	server := new(server.Server)
	log.Println("server start")
	if err := server.ServerRun("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error ocured while running http server: %s", err.Error())
	}
}
