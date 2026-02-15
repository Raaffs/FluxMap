package main

import (
	"context"
	"encoding/json"
	"log"

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
	user := NewWSUser(username, conn, app.websocket)
	app.websocket.AddClient(user)

	go user.KeepAlive(username)
	app.SendUpdateNotification(c.Request().Context(),username)
	app.SendInviteNotification(c.Request().Context(),username)
	return nil
}

func (app *Application)UpdateCountAndNotify(c echo.Context,projectID int, msg, targetUsername string)error{
	username:=c.Get(sessionvar.USERNAME).(string)
	log.Println("Target user name in notify: ",targetUsername)
	if err:=app.models.Update.Create(c.Request().Context(),projectID,msg,username,"user",&targetUsername);err!=nil{
			c.Logger().Error("Error updating updates: ",err)
			return err
	}
	app.SendUpdateNotification(c.Request().Context(),targetUsername)
	return nil
}

func(app *Application)SendUpdateNotification(ctx context.Context, username string){

	count,err:=app.models.Update.GetTotalUnread(ctx,username);if err!=nil{
		app.logger.Error("Error getting total updates: ",err)
	}
	m:=map[NotificationType]int{UpdateNotification:count}
	app.logger.Error("count in update: ",m)
	msg,err:=json.Marshal(m);if err!=nil{
		app.logger.Error("Error marshalling json while sending update notification: ",err)
		return 
	}
	errchan:=app.websocket.PushToClients([]byte(msg),
		func(w WSUser) bool {
		return w.Username==username
	})
	app.CheckChannelError(errchan)
}

func (app *Application)SendInviteNotification(ctx context.Context, username string){
	count,err:=app.models.Invitation.GetTotalUnreadInvitation(ctx,username);if err!=nil{
		app.logger.Error("Error getting invite count: ",err)
		return 
	}

	m:=map[NotificationType]int{InviteNotification:count}

	msg,err:=json.Marshal(m);if err!=nil{
		app.logger.Error("Error marshalling json while sending invite notification: ",err)
		return
	}
	
	errchan:=app.websocket.PushToClients(msg,
		func(w WSUser) bool {
		return w.Username==username
	})

	app.CheckChannelError(errchan)
}

func (app *Application)CheckChannelError( errchan <- chan error){
	if len(errchan)!=0{
		for err:=range errchan{
			app.logger.Warn("error pushing notification: ",err)
		}
	}
}
