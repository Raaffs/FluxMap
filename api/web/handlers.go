package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Raaffs/FluxMap/internal/models"
	validator "github.com/Raaffs/FluxMap/internal/validators"

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
	// v:=validator.New()
	if err:=c.Bind(&u);err!=nil{
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	
	// v.Check(
	// 	validator.MinNameLength(u.Username,3),
	// 	validator.ErrNameTooShort.Key,
	// 	validator.ErrNameTooShort.Message,
	// )

	// v.Check(
	// 	validator.IsStrongPassword(u.Password),
	// 	validator.ErrPasswordTooWeak.Key,
	// 	validator.ErrDescriptionTooShort.Message,
	// )

	hash,err:=HashPassword(u.Password);if err!=nil{
		return echo.NewHTTPError(echo.ErrInternalServerError.Code,"error creating user")
	}

	u.HashedPassword=hash
	if err:=app.models.Users.Create(context.Background(),u);err!=nil{
		if errors.Is(err,models.ErrAlreadyExist){
			return c.JSON(http.StatusConflict,MapMessage("error",models.ErrAlreadyExist.Error()))
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
	// v:=validator.New()
	if err:=c.Bind(&p);err!=nil{
		return echo.NewHTTPError(http.StatusBadRequest,"invalid json")
	}

	// v.Check(
	// 	validator.MinNameLength(p.ProjectName,5),
	// 	validator.ErrNameTooShort.Key,
	// 	fmt.Sprintf(validator.ErrNameTooShort.Message,5),
	// )

	// v.Check(
	// 	validator.MinDescriptionLength(p.ProjectDescription.String,10),
	// 	validator.ErrDescriptionTooShort.Key,
	// 	fmt.Sprintf(validator.ErrDescriptionTooShort.Message,10),
	// )

	// if !v.Valid(){
	// 	return c.JSON(http.StatusBadRequest,v)
	// }

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

func (app *Application)CreateTask(c echo.Context)error{
	var t models.Task
	v:= validator.New()
	fmt.Println("fjiejfoiejfoijeif")
	if err:=c.Bind(&t);err!=nil{
		return c.JSON(http.StatusBadRequest,MapMessage("error",ErrInvalidJson.Error()))
	}
	id,err:=strconv.Atoi(c.Param("id"));if err!=nil{
		c.Logger().Error(MapMessage("error converting to string",err.Error()))
		return c.JSON(http.StatusBadRequest,MapMessage("error","Invalid project ID"))
	}
	t.ParentProjectID=id
	v.Check(
		validator.MinNameLength(t.TaskName,3),
		validator.ErrNameTooShort.Key,
		validator.ErrNameTooShort.Message,
	)
	if !v.Valid(){
		return c.JSON(http.StatusBadRequest,v)
	}
	fmt.Println("task",t)
	if err:=app.models.Task.Create(c.Request().Context(),t);err!=nil{
		c.Logger().Error(MapMessage("Error creating task",err.Error()))
		return c.JSON(http.StatusInternalServerError,MapMessage("error","internal server error"))
	}
	return c.JSON(http.StatusOK,"task created")
}

func (app *Application)GetTasks(c echo.Context)error{
	projectID:=c.Param("id")
	id,err:=strconv.Atoi(projectID);if err!=nil{
		if errors.Is(err,models.ErrRecordNotFound){
			return c.JSON(http.StatusNotFound,MapMessage("error","Project not found"))
		}
		c.Logger().Error(MapMessage("error converting to string",err.Error()))
		return c.JSON(http.StatusNotFound,MapMessage("error","project not found"))
	}
	tasks,err:=app.models.Task.GetTasks(c.Request().Context(),id);if err!=nil{
		c.Logger().Error(MapMessage("Error retrieving tasks",err.Error()))
		return c.JSON(http.StatusInternalServerError,MapMessage("error","internal server error"))
	}
	return c.JSON(http.StatusOK,tasks)
}

func (app *Application)GetTaskByID(c echo.Context)error{
	id,err:=strconv.Atoi(c.Param("taskID"));if err!=nil{
		c.Logger().Error(MapMessage("error converting to string",err.Error()))
		return c.JSON(http.StatusNotFound,MapMessage("error","task not found"))
	}
	task,err:=app.models.Task.GetTaskByID(c.Request().Context(),id);if err!=nil{
		c.Logger().Error(MapMessage("Error retrieving task",err.Error()))
		return c.JSON(http.StatusInternalServerError,MapMessage("error","failed to retrieve task"))
	}

	return c.JSON(http.StatusOK,task)
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
	var t models.Task
	v:=validator.New()
	if err:=c.Bind(&t);err!=nil{
		return c.JSON(http.StatusBadRequest,"Invalid request body")
	}

	v.Check(
		t.TaskDescription.Valid,
		"task",
		"no valid name",
	)

	v.Check(
		t.AssignedUsername.Valid,
		"user",
		"no valid user",
	)

	v.Check(
		t.TaskID>0,
		"id",
		"invalid task id",
	)

	v.Check(
		t.ParentProjectID>0,
		"Projectid",
		"invalid project id",
	)

	v.Check(
		t.Approved.Valid,
		"approved",
		"invalid status",
	)

	if !v.Valid(){
		c.Logger().Error(v)
		c.JSON(http.StatusBadRequest,v)
	}

	if err:=app.models.Task.UpdateManagerTask(c.Request().Context(),t);err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return c.JSON(http.StatusNotFound,map[string]string{"error":"task not found"})
		}
		c.Logger().Error("error updating manager task : ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"failed to update task"})
	}

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
