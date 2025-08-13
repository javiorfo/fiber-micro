package scripts

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, path string) {
	files, err := os.ReadDir(path)
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

		filePath := filepath.Join(path, filename)

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
