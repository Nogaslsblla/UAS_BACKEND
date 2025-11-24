package main

import (
	"log"
	"uas_backend/config"
	"uas_backend/database"
	"uas_backend/route"
)

func main() {
	cfg := config.LoadConfig()
	app := config.NewFiber()

	db := database.ConnectPostgres(cfg)

	route.SetupRoutes(app, db)

	log.Fatal(app.Listen(":" + cfg.AppPort))
}
