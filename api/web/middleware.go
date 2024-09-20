package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func IsAuthorizedUser(next echo.HandlerFunc)echo.HandlerFunc{
	return func (c echo.Context)error{
		cookie,err:=c.Cookie("session"); if err!=nil{
			return c.Redirect(http.StatusTemporaryRedirect,"http://localhost:5173/login")
		}
		if cookie.Value ==""{
			return c.JSON(http.StatusUnauthorized,"Unauthorized")
		}
		return next(c)
	}
}

func IsManager(next echo.HandlerFunc) echo.HandlerFunc {
    return IsAuthorizedUser(func(c echo.Context) error {
        cookie, err := c.Cookie("role")
        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing role cookie"})
        }	

        if cookie.Value != "manager" {
            return c.JSON(http.StatusForbidden, map[string]error{"message": fmt.Errorf("Access denied")})
        }
        return next(c)
    })
}

