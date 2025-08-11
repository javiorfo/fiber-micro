package config

import (
	"os"

	"github.com/javiorfo/fiber-micro/adapter/database/connection"
)

const (
	AppName        = "fiber-micro"
	AppPort        = ":8080"
	AppContextPath = "/app"
)

var DBDataConnection = connection.DBDataConnection{
	Url:         os.Getenv("DATABASE_URL"),
	ShowSQLInfo: true,
}
