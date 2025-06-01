package main

import (
	sessionvar "github.com/Raaffs/FluxMap/internal/sessionVar"
	"github.com/labstack/echo/v4"
)

func (app *Application)HandlWS(c echo.Context)error{
	username:=c.Get(sessionvar.USERNAME).(string)
	app.websocket.KeepAlive(c.Response(),c.Request(),username)
	return nil
}