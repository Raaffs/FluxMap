package main

import (
	"context"
	"fmt"
	"log"
	"github.com/Raaffs/FluxMap/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Application struct{
	env 	map[string]string
	models	models.Models 
}

type config struct {
	port string
	env map[string]string
	db struct {
		dsn string
	}
}


func main(){
	if err:=godotenv.Load(".env");err!=nil{
		log.Fatal("Error  loading .env file %w\n",err)
	}
	envMap,err:=godotenv.Read(".env");if err!=nil{
		log.Fatal("Error reading .env file %w\n",err)
	}
	ctx:=context.Background()

	conn,err:=pgxpool.New(ctx,envMap["DB_URL"]);if err!=nil{
		log.Fatal("Error connecting to database %w\n",err)
	}
	app:=&Application{
		env:	envMap,
		models: models.NewModels(conn),
	}
	e:=app.InitRoutes()
	PORT:=fmt.Sprintf(":%s",app.env["API_PORT"])
	if err:=e.Start(PORT);err!=nil{
		log.Fatal("Error starting server %w\n",err)
	}

}	

