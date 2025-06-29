package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	sessionvar "github.com/Raaffs/FluxMap/internal/sessionVar"
	"github.com/labstack/echo/v4"
)


type NotificationType string

const(
	InviteNotification="inviteNotification"
	UpdateNotification="updateNotification"
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

	app.CheckChannelError(errchan)
	return nil
}

func (app *Application)UpdateCountAndNotify(c echo.Context,projectID int, msg, targetUsername string)error{
	username:=c.Get(sessionvar.USERNAME).(string)
	log.Println("Target user name in notify: ",targetUsername)
	if err:=app.models.Update.CreateUpdate(c.Request().Context(),projectID,msg,username,"user",&targetUsername);err!=nil{
			c.Logger().Error("Error updating updates: ",err)
			return err
	}
	app.SendUpdateNotification(c.Request().Context(),targetUsername)

	return nil
}

func(app *Application)SendUpdateNotification(ctx context.Context, username string){

	count,err:=app.models.Update.GetTotalUnreadUpdates(ctx,username);if err!=nil{
		app.logger.Error("Error getting total updates: ",err)
	}
	m:=map[NotificationType]int{UpdateNotification:count}

	msg,err:=json.Marshal(m);if err!=nil{
		return 
	}
	errchan:=app.websocket.PushToClients([]byte(msg),
		func(w WSUser) bool {
		return w.Username==username
	})
	app.CheckChannelError(errchan)
}

func (app *Application)SendInviteNotification(c echo.Context, username string)error{
	count,err:=app.models.Invitation.GetTotalUnreadInvitation(c.Request().Context(),username);if err!=nil{
		return err
	}

	m:=map[NotificationType]int{InviteNotification:count}

	msg,err:=json.Marshal(m);if err!=nil{
		return err
	}
	
	errchan:=app.websocket.PushToClients(msg,
		func(w WSUser) bool {
		return w.Username==username
	})

	app.CheckChannelError(errchan)
	return nil
}

func (app *Application)CheckChannelError( errchan <- chan error){
	if len(errchan)!=0{
		for err:=range errchan{
			app.logger.Warn("error pushing notification: ",err)
		}
	}
}
