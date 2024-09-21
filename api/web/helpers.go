package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func SetCookie(key string, value string, c echo.Context){
	cookie := &http.Cookie{
        Name:  key,
        Value: value,
        Expires: time.Now().Add(72*time.Hour) ,
    }
	c.SetCookie(cookie)
}

func FormatDate(t time.Time)string{
	return t.Format("dd-mm-yyyy")
}

func HashPassword(password string)(string,error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}