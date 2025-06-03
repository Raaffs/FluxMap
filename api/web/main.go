package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Raaffs/FluxMap/internal/env"
	"github.com/Raaffs/FluxMap/internal/models"
	"github.com/labstack/echo/v4"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Application struct{
	env 			map[string]string
	models			models.Models 
	websocket		*ConnectionManager
}

func main(){
	if err:=godotenv.Load(".env");err!=nil{
		log.Fatal("Error  loading .env file %w\n",err)
	}
	envMap,err:=godotenv.Read(".env");if err!=nil{
		log.Fatal("Error reading .env file %w\n",err)
	}
	ctx:=context.Background()

	conn,err:=pgxpool.New(ctx,envMap[env.DB_URL]);if err!=nil{
		log.Fatal("Error connecting to database %w\n",err)
	}
	
	app:=&Application{
		env:	envMap,
		models: models.NewModels(conn),
		websocket: NewConnectionManager(),
	}
	router:=echo.New()

	app.LoadMiddleware(router)
	app.RegisterRoutes(router)

	// Start the global dispatcher goroutine that listens on ws.Send channel
	// and broadcasts messages to connected clients.
	// This ensures only one goroutine reads from the Send channel,
	// preventing race conditions and message loss.
	PORT:=fmt.Sprintf(":%s",app.env[env.API_PORT])
	if err:=router.Start(PORT);err!=nil{
		log.Fatal("Error starting server %w\n",err)
	}
}	

