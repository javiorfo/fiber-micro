package main

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/javiorfo/fiber-micro/adapter/database/connection"
	"github.com/javiorfo/fiber-micro/adapter/database/migrator/tables"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err == nil {
		log.Info("Using .env file!")
	}

	conn := connection.DBDataConnection{
		Url:         os.Getenv("DATABASE_URL"),
		ShowSQLInfo: false,
	}

	if err := conn.Connect(); err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}
	db := connection.DBinstance

	tables.Migrate(db, "./migrations")
}
