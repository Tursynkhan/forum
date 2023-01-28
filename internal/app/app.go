package app

import (
	"forum/internal/config"
	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
	"log"
)

func Run(cfg *config.Config) {
	db, err := repository.InitDB(repository.Config{
		DBName: cfg.Database.DBName,
		Name:   cfg.Database.Name,
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
	if err := server.ServerRun(cfg, handlers.InitRoutes()); err != nil {
		log.Fatalf("error ocured while running http server:%s", err.Error())
	}
}
