package main

import (
	"context"
	"log"
	"strconv"

	sessionvar "github.com/Raaffs/FluxMap/internal/sessionVar"
	"github.com/labstack/echo/v4"
)

func (app *Application)HandlWS(c echo.Context)error{
	username:=c.Get(sessionvar.USERNAME).(string)
	conn,err:=app.websocket.Upgrader.Upgrade(c.Response(),c.Request(),nil);if err!=nil{
		log.Println("Error creating websocket connection: ",err)
		return nil
	}

	var errchan <-chan error

	user := NewWSUser(username, conn, app.websocket)
	app.websocket.AddClient(user)

	go user.KeepAlive(username)

	count,err:=app.models.Update.GetTotalUnreadUpdates(c.Request().Context(),username);if err!=nil{
		c.Logger().Error("Error getting total unread updates: ",err)
	}

	errchan=app.websocket.PushToClients([]byte(strconv.Itoa(count)),
		func(w WSUser) bool {
			return w.Username==username
	})

	app.CheckChannelError(c,errchan)
	return nil
}


func (app *Application)Send(ctx context.Context, username string){
	if err:=app.models.Update.SetRead(ctx,username);err!=nil{
		log.Println("Error updating read status : ")
	}
}

func (app *Application)SendUpdateNotification(c echo.Context, username string)error{
	count,err:=app.models.Update.GetTotalUnreadUpdates(c.Request().Context(),username);if err!=nil{
		return err
	}

	errchan:=app.websocket.PushToClients([]byte(strconv.Itoa(count)),
		func(w WSUser) bool {
		return w.Username==username
	})

	app.CheckChannelError(c,errchan)
	return nil
}

func (app *Application)SendInviteNotification(c echo.Context, username string)error{
	count,err:=app.models.Invitation.GetTotalUnreadInvitation(c.Request().Context(),username);if err!=nil{
		return err
	}

	errchan:=app.websocket.PushToClients([]byte(strconv.Itoa(count)),
		func(w WSUser) bool {
		return w.Username==username
	})

	app.CheckChannelError(c,errchan)
	return nil
}


func (app *Application)CheckChannelError(c echo.Context, errchan <- chan error){
	if len(errchan)!=0{
		for err:=range errchan{
			c.Logger().Warn("error pushing notification: ",err)
		}
	}
}
