package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SetCookie(key string, value string, c echo.Context){
	cookie := &http.Cookie{
        Name:  key,
        Value: value,
        Expires: <-time.After(72*time.Hour),
    }
	c.SetCookie(cookie)
}