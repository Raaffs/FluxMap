package main

import (
)

type Application struct{
	users string	
}

func main(){
	app:=&Application{}
	e:=app.InitRoutes()
	e.Start(":8080")
}	

