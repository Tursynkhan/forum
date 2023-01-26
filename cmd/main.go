package main

import (
	"forum/internal/app"
	"forum/internal/config"
	"log"
)

func main() {
	cfg, err := config.InitConfig("./internal/config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	app.Run(cfg)
}
