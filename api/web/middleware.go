package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func IsAuthorizedUser(next echo.HandlerFunc)echo.HandlerFunc{
	return func (c echo.Context)error{
		cookie,err:=c.Cookie("session"); if err!=nil{
			return c.Redirect(http.StatusTemporaryRedirect,"http://localhost:5173/login")
		}
		if cookie.Value == ""{
			return c.JSON(http.StatusUnauthorized,"Unauthorized")
		}
		return next(c)
	}
}

func(app *Application)ManagerLevelAccess(next echo.HandlerFunc) echo.HandlerFunc {
    return IsAuthorizedUser(func(c echo.Context) error {
        cookie, err := c.Cookie("username")
        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing role cookie"})
        }

		isAdmin,err:=app.models.Users.IsAdmin(c.Request().Context(),cookie.Value,c.Param("id"));if err!=nil{
			return c.JSON(http.StatusInternalServerError,map[string]string{"message":err.Error()})
		}
		if isAdmin{
			return next(c)
		}

		isManager,err:=app.models.Users.IsManager(c.Request().Context(),cookie.Value,c.Param("id"))
		if err!=nil{
			return c.JSON(http.StatusInternalServerError,map[string]string{"message":err.Error()})
		}
		if !isManager{
			return c.JSON(http.StatusForbidden,map[string]string{"message":"You are not a manager"})
		}
		return next(c)
    })
}

func (app *Application)AdminLevelAccess(next echo.HandlerFunc)echo.HandlerFunc{
	return IsAuthorizedUser(func(c echo.Context) error {
		cookie, err := c.Cookie("username")
        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing role cookie"})
        }

		isAdmin,err:=app.models.Users.IsAdmin(c.Request().Context(),cookie.Value,c.Param("id"));if err!=nil{
			return c.JSON(http.StatusInternalServerError,map[string]string{"message":err.Error()})
		}
		if !isAdmin{
			return c.JSON(http.StatusForbidden,map[string]string{"message":"You are not an admin"})
		}
		return next(c)
	})
}

