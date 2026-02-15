package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Raaffs/FluxMap/internal/env"
	"github.com/Raaffs/FluxMap/internal/graphs"
	"github.com/Raaffs/FluxMap/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)	

type Application struct {
	env       	 map[string]string
	models    	 models.Models
	graphs    	 graphs.GraphModel
	websocket 	*ConnectionManager
	logger    	 echo.Logger
	oAuthConfig *oauth2.Config
}


func connectWithRetry(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	for i := range 10 {
		log.Printf("Attempting DB connection (attempt %d/10)...", i+1)
		
		pool, err = pgxpool.New(ctx, dbURL)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				log.Println("PING SUCCESSFUL")
				return pool, nil
			}
		}

		log.Printf("DB not ready: %v. Retrying in 2s...", err)
		if pool != nil {
			pool.Close()
		}
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to DB: %w", err)
}

func main() {

	envMap := map[string]string{
		env.API_PORT:      	os.Getenv(env.API_PORT),
		env.DB_URL:        	os.Getenv(env.DB_URL),
		env.CLIENT_PORT:   	os.Getenv(env.CLIENT_PORT),
		env.SESSION_SECRET:	os.Getenv(env.SESSION_SECRET),
		env.GOOGLE_OAUTH_CLIENT_ID: os.Getenv(env.GOOGLE_OAUTH_CLIENT_ID),
		env.GOOGLE_OAUTH_SECRET: os.Getenv(env.GOOGLE_OAUTH_SECRET),
	}

	oauthConfig := &oauth2.Config{
		ClientID:     envMap[env.GOOGLE_OAUTH_CLIENT_ID],
		ClientSecret: envMap[env.GOOGLE_OAUTH_SECRET],
		RedirectURL:  "http://localhost:4000/api/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
	log.Println(oauthConfig.ClientID,oauthConfig.ClientSecret)
	if envMap[env.DB_URL] == "" {
		log.Fatal("DB_URL is missing")
	}

	ctx := context.Background()
	
	conn, err := connectWithRetry(ctx,envMap[env.DB_URL]);if err!=nil{
		log.Fatalf("Could not connect to DB: %v", err)
	}

	app := &Application{
		env:       envMap,
		models:    models.NewModels(conn),
		graphs:    graphs.NewGraphs(conn),
		websocket: NewConnectionManager(),
		logger:    echo.New().Logger,
		oAuthConfig: oauthConfig,
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

