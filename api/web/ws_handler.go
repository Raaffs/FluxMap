package main

import (
	"fmt"
	"log"

	sessionvar "github.com/Raaffs/FluxMap/internal/sessionVar"
	"github.com/labstack/echo/v4"
)

func (app *Application)HandlWS(c echo.Context)error{
	username:=c.Get(sessionvar.USERNAME).(string)
	conn,err:=app.websocket.Upgrader.Upgrade(c.Response(),c.Request(),nil);if err!=nil{
		log.Println("Error creating websocket connection: ",err)
		return nil
	}
	errchan:=make(<-chan error)

	user := NewWSUser(username, conn, app.websocket)
	app.websocket.AddClient(user)

	go user.KeepAlive(username)
	//example of push to client:

		send:=fmt.Sprintf("%s joined the server",username)
		errchan=app.websocket.PushToClients([]byte(send),func(w WSUser) bool {return true})
	
	for err:=range errchan{
		log.Print(err)
	}
	return nil
}