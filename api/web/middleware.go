package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

        username:=c.Get(sessionvar.USERNAME).(string)
		// We're manually reading the request body instead of using c.Bind() because:
		// 1. c.Bind() consumes the body, meaning the next handler (Create/Update Task) wouldn't be able to read it.
		// 2. To avoid that, we read the body into a variable, then reset c.Request().Body so it can be read again.
		// 3. This makes sure everything works smoothly without breaking the request flow.
        bodyBytes, err := io.ReadAll(c.Request().Body)
        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to read request body"})
        }
        c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

        Notify := struct {
            TargetUsername string `json:"targetUsername"`
            TargetType     string `json:"targetType"`
            Msg            string `json:"msg"`
        }{}

        if len(bodyBytes) > 0 {
            if err := json.Unmarshal(bodyBytes, &Notify); err != nil {
                c.Logger().Error("Error sending notification\nerror binding json: ", err)
                return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
            }
        }
        log.Println("NOTIFYING THE USER: ",Notify)
        id := c.Param("id")
        projectID, err := strconv.Atoi(id)
        if err != nil {
            c.Logger().Error("Error sending notification\ninvalid projectid: ", err)
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
        }

        if Notify.TargetType != "user" && Notify.TargetType != "all" {
            c.Logger().Error("Error sending notification\ninvalid target type: ", Notify.TargetType)
        }



        role, err := app.getUserRole(c.Request().Context(), username, id)
        if err != nil {
            c.Logger().Error("Error sending notification\nerror getting user role: ", err)
        }

        var msg string
        if Notify.Msg == "" {
            switch role {
            case AdminRole:
                msg = fmt.Sprintf("Task %s has been approved by %s", id, username)
            case ManagerRole:
                msg = fmt.Sprintf("Task %s has been approved by %s", id, username)
            case UserRoleVal:
                msg = fmt.Sprintf("Task %s has been completed by %s", id, username)
            default:
                c.Logger().Error("Unexpected role, message is empty")
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unexpected role"})
            }
        } else {
            msg = fmt.Sprintf("%s %s",Notify.Msg,username)
        }

        // Defer notification, so it's only sent if the handler executes successfully
        defer func() {
            if err!=nil{
                c.Logger().Error("Error sending notification: ", err, Notify)
                return
            }
            if c.Response().Status >= 200 && c.Response().Status < 300 { // Only proceed if the handler succeeds
                if err := app.models.Notify.CreateUpdate(c.Request().Context(), projectID, msg, username, Notify.TargetType, Notify.TargetUsername); err != nil {
                    c.Logger().Error("Error sending notification\nerror creating notification: ", err)
                }
            }
        }()

        return next(c)
    }
}   