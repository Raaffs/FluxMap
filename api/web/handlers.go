package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
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
		c.Logger().Error("Error authenticating user: ",err)
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
		c.Logger().Error("Error creating user: ",err)
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
	if err:=c.Bind(&p);err!=nil{
		return echo.NewHTTPError(http.StatusBadRequest,"invalid json")
	}
	if err:=app.models.Projects.Create(c.Request().Context(),p);err!=nil{
		c.Logger().Error("Error creating project: ",err)
		return c.JSON(http.StatusInternalServerError,"An error while creating project")
	}
	return c.JSON(http.StatusOK,"project created")
}

func (app *Application) GetProjects(c echo.Context) error {
	user, err := c.Cookie("username")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	username := user.Value

	projects := struct {
		AdminProjects    []*models.Project `json:"adminProjects"`
		ManagerProjects  []*models.Project `json:"managerProjects"`
		AssignedProjects []*models.Project `json:"assignedProjects"`
	}{}

	adminchan := make(chan []*models.Project)
	managerchan := make(chan []*models.Project)
	assginedchan := make(chan []*models.Project)
	errorchan := make(chan error)
	done:=make(chan struct{})
	log.Println("username:", username)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()
	
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		projects, err := app.models.Projects.RetrieveAdminProjects(ctx, username)
		if err != nil {
			errorchan <- err
			return
		}
		adminchan <- projects
	}()

	go func() {
		defer wg.Done()
		projects, err := app.models.Projects.RetrieveManagerProjects(ctx, username)
		if err != nil {
			errorchan <- err
			return
		}
		managerchan <- projects
	}()

	go func() {
		defer wg.Done()
		projects, err := app.models.Projects.RetrieveAssginedProjects(ctx, username)
		if err != nil {
			c.Logger().Error("Error retrieving projects: ",err)
			errorchan <- err
			return
		}
		assginedchan <- projects
	}()


	go func() {
		wg.Wait()
		close(adminchan)
		close(managerchan)
		close(assginedchan)
		close(errorchan)
		close(done)
	}()

	for {
		select {
		case admin, ok := <-adminchan:
			if ok {
				projects.AdminProjects = admin
			}
		case manager, ok := <-managerchan:
			if ok {
				projects.ManagerProjects = manager
			}
		case assgined,ok:=<-assginedchan:
			if ok{
				projects.AssignedProjects=assgined
			}
		case err := <-errorchan:
			if err!=nil{
				c.Logger().Error("Error retrieving projects: ", err)
				return c.JSON(http.StatusInternalServerError, "An error while retrieving projects")	
			}
		case <-done:
			return c.JSON(http.StatusOK, projects)
		}
	}
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

func(app *Application)AddManager(c echo.Context)error{
	m:=struct{
		Manager string `json:"manager"`
	}{}
	fmt.Println("herer")
	projectID:=c.Param("id")
	fmt.Println(m.Manager,projectID)
	fmt.Println(projectID,m.Manager)
	if err:=c.Bind(&m);err!=nil{
		return c.JSON(http.StatusBadRequest,"Invalid request body")
	}
	if err:=app.models.Projects.AssignManager(c.Request().Context(),m.Manager,projectID);err!=nil{
		return c.JSON(http.StatusInternalServerError,"Failed to assign manager")
	}
	return c.JSON(http.StatusOK,"manager added")
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
