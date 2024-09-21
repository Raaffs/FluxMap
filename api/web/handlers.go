package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Raaffs/FluxMap/internal/models"

	"github.com/labstack/echo/v4"
)

func(app *Application)Home(c echo.Context)error{
	data:=map[string]any{
		"message":"Hello,World!",
		"id":2,
	}
	if data["id"]==2{
		SetCookie("role","manager",c)
		return c.JSON(http.StatusOK,data)
	}
	return c.JSON(http.StatusBadRequest,"not found")
}

func(app *Application)Login(c echo.Context)error{
	var u models.User
	err := c.Bind(&u); if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

	if err:=app.models.Users.Login(c.Request().Context(),u.Username,u.Password);err!=nil{
		if errors.Is(err,models.ErrInvalidCredential){
			return c.JSON(http.StatusUnauthorized,"invalid credential")
		}
									
		if errors.Is(err,sql.ErrNoRows){
			return c.JSON(http.StatusNotFound,"user not found")
		}
		return c.JSON(http.StatusInternalServerError,err.Error())
	}

	SetCookie("username",u.Username,c)
	return c.JSON(http.StatusOK,"")
}


func(app *Application)Register(c echo.Context)error{
	u:=models.User{}
	
	if err:=c.Bind(&u);err!=nil{
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	
	hash,err:=HashPassword(u.Password);if err!=nil{
		return echo.NewHTTPError(echo.ErrInternalServerError.Code,"error creating user")
	}

	u.HashedPassword=hash
	if err:=app.models.Users.Create(context.Background(),u);err!=nil{
		if errors.Is(err,models.ErrAlreadyExist){
			return c.JSON(http.StatusConflict,"User already exist")
		}
		return echo.NewHTTPError(echo.ErrInternalServerError.Code,"internal server error")
	}

	SetCookie("username",u.Username,c)
	return c.JSON(http.StatusOK,"user registered successfully")
}

func(app *Application)Logout(c echo.Context)error{
	cookie := &http.Cookie{
		Name:     "role",
		Value:    "",
		Path:     "/",
		Expires: time.Unix(0, 0),
		HttpOnly: true,
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK,"")
}

	func (app *Application)CreateProject(c echo.Context)error{
		var p models.Project
		// body, err := io.ReadAll(c.Request().Body)
		// if err != nil {
		//     return echo.NewHTTPError(http.StatusBadRequest, "unable to read request body")
		// }
		// if err := json.Unmarshal(body, &p); err != nil {
		//     return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
		// }
		if err:=c.Bind(&p);err!=nil{
			return echo.NewHTTPError(http.StatusBadRequest,"invalid json")
		}
		if err:=app.models.Projects.Create(c.Request().Context(),p);err!=nil{
			return c.JSON(http.StatusInternalServerError,err.Error())
		}
		return c.JSON(http.StatusOK,"project created")
	}

func (app *Application)GetProjects(c echo.Context)error{
	return c.JSON(http.StatusOK,"projects retrieved")
}

func (app *Application)GetProjectByID(c echo.Context)error{
	return c.JSON(http.StatusOK,"project retrieved")
}

func (app *Application)CreateTask(c echo.Context)error{
	return c.JSON(http.StatusOK,"task created")
}

func (app *Application)GetTasks(c echo.Context)error{
	return c.JSON(http.StatusOK,"tasks retrieved")
}

func (app *Application)GetTaskByID(c echo.Context)error{
	return c.JSON(http.StatusOK,"task retrieved")
}

func (app *Application)GetAssignedUserByProject(c echo.Context)error{
	return c.JSON(http.StatusOK,"assgined users retrived")
}

func (app *Application)GetAssignedUserByTask(c echo.Context)error{
	return c.JSON(http.StatusOK,"assgined user by task retrived")
}

func (app *Application)UpdateProject(c echo.Context)error{
	return c.JSON(http.StatusOK,"updated successfully")
}

func (app *Application)UpdateTask(c echo.Context)error{
	return c.JSON(http.StatusOK,"updated successfully")
}

func (app *Application)ManagerRestrictedTask(c echo.Context)error{
	return c.JSON(http.StatusOK,"task approved")
}

func (app *Application)AdminRestrictedProject(c echo.Context)error{
	return c.JSON(http.StatusOK,"task approved")
}


func(app *Application)GetPert(c echo.Context)error{
	return c.JSON(http.StatusOK,"done")
}

func(app *Application)CreatePert(c echo.Context)error{
	return c.JSON(http.StatusOK,"done")
}

func(app *Application)UpdatePert(c echo.Context)error{
	return c.JSON(http.StatusOK,"done")
}

func(app *Application)GetCpm(c echo.Context)error{
	return c.JSON(http.StatusOK,"done")
}

func(app *Application)CreateCpm(c echo.Context)error{
	return c.JSON(http.StatusOK,"done")
}

func(app *Application)UpdateCpm(c echo.Context)error{
	return c.JSON(http.StatusOK,"done")
}
