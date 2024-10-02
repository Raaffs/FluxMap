package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var(
    ErrInvalidJson=errors.New("Invalid JSON")
    ErrFetchingResult=errors.New("Error getting analytics")
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

func MapMessage(key string,msg string)struct{Message map[string]string}{
    return struct{Message map[string]string}{
        Message: map[string]string{
            key:key,
            msg:msg,
        },
    }
}

func HashPassword(password string)(string,error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

