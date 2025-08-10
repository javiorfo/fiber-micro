package main

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/gofiber/fiber/v2/log"
	"github.com/javiorfo/fiber-micro/adapter/database/connection"

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

	scriptsPath := "./deploy/scripts"

	files, err := os.ReadDir(scriptsPath)
	if err != nil {
		log.Fatalf("failed to read migrations directory: %v", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)

	for _, filename := range migrationFiles {
		log.Infof("Executing migration: %s", filename)

		filePath := filepath.Join(scriptsPath, filename)

		sqlScript, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("failed to read migration file %s: %v", filePath, err)
		}

		if err := db.Exec(string(sqlScript)).Error; err != nil {
			log.Fatalf("failed to execute migration %s: %v", filePath, err)
		}
	}

	log.Info("Migration completed succesfully!")
}
