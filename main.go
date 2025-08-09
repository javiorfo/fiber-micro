package main

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err == nil {
		log.Info("Using .env file!")
	}
}
