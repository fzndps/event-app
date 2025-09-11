package main

import (
	"log"
	"rest-api-event-app/cmd/migrate"
	_ "rest-api-event-app/docs"
	"rest-api-event-app/internal/database"
	"rest-api-event-app/internal/env"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

// @title           Rest API Event App
// @version         1.0
// @description     Creating a REST API event app with Gin and JWT
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConn, err := migrate.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s", err)
	}

	defer dbConn.CloseDB()

	models := database.NewModels(dbConn.GetDB())
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-150902"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
