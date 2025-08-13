package main

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/javiorfo/fiber-micro/adapter/database"
	"github.com/javiorfo/fiber-micro/adapter/database/migrator/scripts"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err == nil {
		log.Info("Using .env file!")
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}
	db := database.DBinstance

	scripts.Migrate(db, "./migrations")
}
