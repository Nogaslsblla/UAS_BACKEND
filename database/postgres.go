package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"uas_backend/config"
)

func ConnectPostgres(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open DB: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB: ", err)
	}

	log.Println("Connected to PostgreSQL")
	return db
}
