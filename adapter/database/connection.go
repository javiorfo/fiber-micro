package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBinstance *gorm.DB

func Connect() error {
	dsn := fmt.Sprintf("%s?sslmode=disable", os.Getenv("DATABASE_URL"))

	isShowSQLInfo, err := strconv.ParseBool(os.Getenv("SHOW_SQL_INFO"))
	if err != nil {
		return fmt.Errorf("Could not parse SHOW_SQL_INFO var %v", err)
	}

	loggerSQL := logger.Default.LogMode(logger.Info)
	if !isShowSQLInfo {
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
