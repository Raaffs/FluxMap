package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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
        HttpOnly: true,
        Secure: false,
        SameSite: http.SameSiteDefaultMode  ,
        Expires: time.Now().Add(72*time.Hour) ,
    }
	c.SetCookie(cookie)
}

func FormatDate(t time.Time)string{
	return t.Format("dd-mm-yyyy")
}


func MapMessage(key string,msg string)struct{Key string; Message string}{
    return struct{
        Key string
        Message string
    }{
        Key: key,
        Message: msg,
    }
}


func AppendToSessionArray(session *sessions.Session, key string, value int) error {
	// Try to assert the session value to []string
	arr, ok := session.Values[key].([]int)
	if !ok {
		// Return an error if the type assertion fails
		return errors.New("session value is not of the expected type []string")
	}

	// Append the new value to the existing array
	session.Values[key] = append(arr, value)

	return nil
}


func addToSession(c echo.Context, key string, value int)error{
    sess,err:=session.Get("session",c); if err!=nil{
        return err
    }
    sess.Values[key] = value
    
    if err:=sess.Save(c.Request(),c.Response().Writer);err!=nil{
        return err
    }

    return nil
}

func HashPassword(password string)(string,error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
       return "", err
    }
    return string(hashedPassword), nil
}

