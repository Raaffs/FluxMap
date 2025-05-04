package main

import (
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

func(app *Application)ManagerLevelAccess(next echo.HandlerFunc) echo.HandlerFunc {
    return IsAuthorizedUser(func(c echo.Context) error {
		
        username:=c.Get(sessionvar.USERNAME).(string)
		isAdmin,err:=app.models.Users.IsAdmin(c.Request().Context(),username,c.Param("id"));if err!=nil{
			c.Logger().Error("error getting access level: ",err)
			return c.JSON(http.StatusInternalServerError,map[string]string{"message":err.Error()})
		}

		if isAdmin{
            c.Set("isAdmin",true)
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

func (app *Application) SendNotification(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        id:=c.Param("id")
        projectID,err:=strconv.Atoi(id);if err!=nil{
            c.Logger().Error("Error sending notification, invalid projectid : ",err)
        }
        // Defer notification, so it's only sent if the handler executes successfully
        defer func() {
            var targetUsername *string

            username:=c.Get(sessionvar.USERNAME).(string)
            msg,ok:=c.Get(NOTIFY_MSG).(string);if !ok {c.Logger().Error("error converting notify msg to string")}
            targetType,ok:=c.Get(NOTIFY_TARGET_TYPE).(string);if !ok {c.Logger().Error("error converting notify target type to string")}
            
            if targetType=="user"{
                tu,ok:=c.Get(NOTIFY_TARGET_USERNAME).(string);if ok {
                    targetUsername=&tu
                }
            }
            log.Println(username,msg,targetType,targetUsername)
            if err!=nil{
                c.Logger().Error("Error sending notification: ", err)
                return
            }
            if c.Response().Status >= 200 && c.Response().Status < 300 { // Only proceed if the handler succeeds
                if err := app.models.Notify.CreateUpdate(c.Request().Context(), projectID, msg, username, targetType, targetUsername); err != nil {
                    c.Logger().Error("Error sending notification\nerror creating notification: ", err)
                }
            }
        }()
        return next(c)
    }
}   