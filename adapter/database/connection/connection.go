package connection

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBDataConnection struct {
	Url         string
	ShowSQLInfo bool
}

var DBinstance *gorm.DB

func (db DBDataConnection) Connect() error {
	dsn := fmt.Sprintf("%s?sslmode=disable", db.Url)

	loggerSQL := logger.Default.LogMode(logger.Info)
	if !db.ShowSQLInfo {
		loggerSQL = logger.Discard
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: loggerSQL,
	})
	if err != nil {
		return err
	}

	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("Could not get sql DB %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(20 * time.Minute)

	err = database.Use(otelgorm.NewPlugin())
	if err != nil {
		return fmt.Errorf("Could not set otelgorm %v", err)
	}

	log.Info("Connected to DB!")
	database.Logger = loggerSQL
	DBinstance = database

	return nil
}
