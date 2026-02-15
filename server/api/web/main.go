package main

import (
	"context"
	"encoding/json"
	"flag"
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

type Config struct {
	Port  int `json:"port"`
	OAuth struct {
		OAuthRedirectUri   		string   `json:"oAuthRedirectUri"`
		NewAccountRedirectUri  	string   `json:"newAccountRedirectUri"`
		SuccessRedirectUri 		string   `json:"successRedirectUri"`
		Scopes             		[]string `json:"scopes"`
	} `json:"oAuth"`
	TrustedOrigins []string `json:"trustedOrigins"`
}

func LoadConfig(path string, env string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var allConfigs map[string]Config
	if err := json.Unmarshal(file, &allConfigs); err != nil {
		return nil, err
	}

	config, ok := allConfigs[env]
	if !ok {
		return nil, fmt.Errorf("environment '%s' not found in config", env)
	}

	return &config, nil
}

type Application struct {
	env       	 map[string]string
	config   	 *Config
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
	c := flag.String("env", "dev", "target environment")
	
	flag.Parse()
	cfg, err := LoadConfig("config.json", *c)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	log.Printf("Loaded config for environment: %s\n", cfg)
	envMap := map[string]string{
		env.API_PORT:      	os.Getenv(env.API_PORT),
		env.DB_URL:        	os.Getenv(env.DB_URL),
		env.CLIENT_PORT:   	os.Getenv(env.CLIENT_PORT),
		env.SESSION_SECRET:	os.Getenv(env.SESSION_SECRET),
		env.GOOGLE_OAUTH_CLIENT_ID: os.Getenv(env.GOOGLE_OAUTH_CLIENT_ID),
		env.GOOGLE_OAUTH_SECRET: os.Getenv(env.GOOGLE_OAUTH_SECRET),
	}
	if envMap[env.DB_URL] == "" {
		log.Fatal("DB_URL is missing")
	}

	oauthConfig := &oauth2.Config{
		ClientID:     envMap[env.GOOGLE_OAUTH_CLIENT_ID],
		ClientSecret: envMap[env.GOOGLE_OAUTH_SECRET],
		RedirectURL:  cfg.OAuth.OAuthRedirectUri,
		Scopes:       cfg.OAuth.Scopes,
		Endpoint: google.Endpoint,
	}
	ctx := context.Background()
	
	conn, err := connectWithRetry(ctx,envMap[env.DB_URL]);if err!=nil{
		log.Fatalf("Could not connect to DB: %v", err)
	}



	app := &Application{
		env:       envMap,
		config:    cfg,
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

