package main

import (
	"log"
	"uas_backend/config"
	"uas_backend/database"
	"uas_backend/route"

	_ "uas_backend/docs" // Import swagger docs

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// Alias untuk swagger wrapper
var swaggerWrapper = fiberSwagger.WrapHandler

// @title UAS Backend API
// @version 1.0
// @description API untuk sistem manajemen prestasi mahasiswa
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.LoadConfig()
	app := config.NewFiber()

	db := database.ConnectPostgres(cfg)
	mongoClient := database.ConnectMongoDB(cfg.MongoURI)
	mongoColl := database.GetCollection(mongoClient, cfg.MongoDBName, "achievements")

	route.SetupRoutes(app, db, mongoColl)

	// Swagger route
	app.Get("/swagger/*", swaggerWrapper)

	log.Printf("Server running on port %s", cfg.AppPort)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", cfg.AppPort)

	log.Fatal(app.Listen(":" + cfg.AppPort))
}
