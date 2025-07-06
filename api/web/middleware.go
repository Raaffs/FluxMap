package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Raaffs/FluxMap/internal/sessionVar"
	"github.com/labstack/echo/v4"
)

func IsAuthorizedUser(next echo.HandlerFunc)echo.HandlerFunc{
	return func (c echo.Context)error{
		username,err:=GetUsernameFromSession(c); if err!=nil{
            c.Logger().Warn("error getting session:", err)
            return c.JSON(http.StatusUnauthorized, map[string]string{"error":"unauthorized"})
        }
        log.Println("username : ",username,c.Path())
        c.Set(sessionvar.USERNAME,username)
		return next(c)
	}
}

func (app *Application)HasProjectAccess(next echo.HandlerFunc)echo.HandlerFunc{
	return func(c echo.Context) error {
		username,ok:=c.Get(sessionvar.USERNAME).(string);if!ok{
			return c.JSON(http.StatusUnauthorized,map[string]string{"error":"unauthorized access"})
		}
		projectID,err:=strconv.Atoi(c.Param("id"));if err!=nil{
			c.Logger().Error("error converting project id string to int: ",err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "project not found"})
		}
		isAdmin,err:=app.models.Users.IsAdmin(c.Request().Context(),username,strconv.Itoa(projectID));if err!=nil{
			if errors.Is(err,sql.ErrNoRows){
				c.Logger().Error("error getting access level: ",err)
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized access"})
			}
			c.Logger().Error("error getting access level: ",err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
		if isAdmin{
			return next(c)
		}

		switch status:=app.EnsureExists(c,func(ctx context.Context) (bool, error) {
		return app.models.Invitation.HasAcceptedInvitation(ctx,projectID,username)
		});status{
		case ErrorCheckingExistStatus:
		// error response already sent by EnsureExists, just exit
			return nil
		case Exists: 
			return next(c)
		case NotExists:
			    c.Logger().Infof("user %s has no accepted invitation to project %d", username, projectID)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized access"})
		default:
			return c.JSON(http.StatusInternalServerError,"If this happens I'm quitting")
		}
	}
}

func(app *Application)ManagerLevelAccess(next echo.HandlerFunc) echo.HandlerFunc {
    return IsAuthorizedUser(func(c echo.Context) error {
        username:=c.Get(sessionvar.USERNAME).(string)
		isAdmin,err:=app.models.Users.IsAdmin(c.Request().Context(),username,c.Param("id"));if err!=nil{
			c.Logger().Error("error getting access level: ",err)
			return c.JSON(http.StatusInternalServerError,map[string]string{"message":err.Error()})
		}

		if isAdmin{
			return next(c)
		}

		isManager,err:=app.models.Users.IsManager(c.Request().Context(),username,c.Param("id"))
		if err!=nil{
			c.Logger().Error("error getting access level: ",err)
			return c.JSON(http.StatusInternalServerError,map[string]string{"message":err.Error()})
		}
		if !isManager{
			c.Logger().Warn("unauthorized access\n","session: ","\nerror: ",err)
			return c.JSON(http.StatusForbidden,map[string]string{"message":"You are not a manager"})
		}

		return next(c)
    })
}

func (app *Application)AdminLevelAccess(next echo.HandlerFunc)echo.HandlerFunc{
	return IsAuthorizedUser(func(c echo.Context) error {
        username:=c.Get(sessionvar.USERNAME).(string)
		isAdmin,err:=app.models.Users.IsAdmin(c.Request().Context(),username,c.Param("id"));if err!=nil{
			return c.JSON(http.StatusInternalServerError,map[string]string{"message":err.Error()})
		}
		if !isAdmin{
		return c.JSON(http.StatusForbidden,map[string]string{"message":"You are not an admin"})
		}
		return next(c)
	})
}
