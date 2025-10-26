package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/swagger"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/javiorfo/fiber-micro/adapter/database"
	_ "github.com/javiorfo/fiber-micro/docs"
	"github.com/javiorfo/fiber-micro/internal"
	"github.com/javiorfo/go-microservice-lib/tracing"
	"github.com/joho/godotenv"
)

const (
	appName        = "fiber-micro"
	appPort        = ":8080"
	appContextPath = "/app"
)

// @contact.name				API Support
// @contact.email				fiber@swagger.io
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath					/app
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	// Env variables
	err := godotenv.Load()
	if err == nil {
		log.Info("Using .env file!")
	}

	// Database
	err = database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database.\n", err)
	}

	// Tracing
	traceProvider, err := tracing.StartTracing(appName)
	if err != nil {
		log.Fatalf("traceprovider: %v", err)
	}
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("traceprovider: %v", err)
		}
	}()

	_ = traceProvider.Tracer(appName)

	// Fiber
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		IdleTimeout:           15 * time.Second,
	})

	app.Use(cors.New())
	app.Use(recover.New())

	swaggerPath := fmt.Sprintf("%s/swagger", appContextPath)
	app.Use(otelfiber.Middleware(otelfiber.WithNext(func(ctx *fiber.Ctx) bool {
		return strings.HasPrefix(ctx.Path(), swaggerPath)
	})))
	log.Info("Tracing configured!")

	// Logger
	err = os.MkdirAll("var/log", 0755)
	if err != nil {
		log.Fatalf("error creating directory var/log: %v", err)
	}

	file, err := os.OpenFile(fmt.Sprintf("./var/log/%s.log", appName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	iw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(iw)
	defer file.Close()

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${method} ${ip}${path} | response: ${status} ${error} | ${latency}\n",
		TimeFormat: "2006/01/02 15:04:05.000000",
		Output:     iw,
		Next:       func(c *fiber.Ctx) bool { return strings.HasPrefix(c.Path(), swaggerPath) },
	}))

	api := app.Group(appContextPath)

	internal.InjectDependencies(api)

	// Swagger
	if os.Getenv("SWAGGER_ENABLED") == "true" {
		app.Get(swaggerPath+"/*", swagger.New(swagger.Config{
			DeepLinking:  false,
			DocExpansion: "list",
		}))
	}

	log.Infof("Context path: %s", appContextPath)
	log.Infof("Starting %s on port %s...", appName, appPort)
	log.Info("Server Up!")

	go func() {
		log.Panic(app.Listen(appPort))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c
	log.Info("Gracefully shutting down...")
	_ = app.Shutdown()
	db, _ := database.DBinstance.DB()
	db.Close()
	log.Info("Done!")
}
