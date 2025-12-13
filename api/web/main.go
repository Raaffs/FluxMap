package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Raaffs/FluxMap/internal/env"
	"github.com/Raaffs/FluxMap/internal/graphs"
	"github.com/Raaffs/FluxMap/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type Application struct {
	env       map[string]string
	models    models.Models
	graphs    graphs.GraphModel
	websocket *ConnectionManager
	logger    echo.Logger
}

func main() {

	// Load .env ONLY if it exists (local dev)
	_ = godotenv.Load()

	// Copy env map EXACTLY as .env defines it
	envMap := map[string]string{
		env.API_PORT:      os.Getenv(env.API_PORT),
		env.DB_URL:        os.Getenv(env.DB_URL),
		env.CLIENT_PORT:   os.Getenv(env.CLIENT_PORT),
		env.SESSION_SECRET: os.Getenv(env.SESSION_SECRET),
	}

	// Hard fails for required vars
	if envMap[env.DB_URL] == "" {
		log.Fatal("DB_URL is missing")
	}
	if envMap[env.API_PORT] == "" {
		log.Fatal("API_PORT is missing")
	}

	ctx := context.Background()

	conn, err := pgxpool.New(ctx, envMap[env.DB_URL])
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	app := &Application{
		env:       envMap,
		models:    models.NewModels(conn),
		graphs:    graphs.NewGraphs(conn),
		websocket: NewConnectionManager(),
		logger:    echo.New().Logger,
	}

	router := echo.New()

	app.LoadMiddleware(router)
	app.RegisterRoutes(router)

	port := fmt.Sprintf(":%s", envMap[env.API_PORT])
	log.Println("API starting on", port)

	if err := router.Start(port); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
