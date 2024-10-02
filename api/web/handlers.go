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

	"github.com/Raaffs/FluxMap/api/external"
	"github.com/Raaffs/FluxMap/internal/models"
	validator "github.com/Raaffs/FluxMap/internal/validators"

	"github.com/labstack/echo/v4"
)

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
	v:=validator.New()
	if err:=c.Bind(&u);err!=nil{
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	
	v.Check(
		validator.MinNameLength(u.Username,3),
		validator.ErrNameTooShort.Key,
		validator.ErrNameTooShort.Message,
	)

	v.Check(
		validator.IsStrongPassword(u.Password),
		validator.ErrPasswordTooWeak.Key,
		validator.ErrDescriptionTooShort.Message,
	)

	v.Check(
		validator.Matches(u.Email,validator.EmailRX),
		"email",
		"invalid email",
	)

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
		Name:     "username",
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
	v:=validator.New()
	if err:=c.Bind(&p);err!=nil{
		return echo.NewHTTPError(http.StatusBadRequest,"invalid json")
	}

	v.Check(
		validator.MinNameLength(p.ProjectName,5),
		validator.ErrNameTooShort.Key,
		fmt.Sprintf(validator.ErrNameTooShort.Message,5),
	)

	v.Check(
		validator.MinDescriptionLength(p.ProjectDescription.String,10),
		validator.ErrDescriptionTooShort.Key,
		fmt.Sprintf(validator.ErrDescriptionTooShort.Message,10),
	)

	if !v.Valid(){
		return c.JSON(http.StatusBadRequest,v)
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
		case <-ctx.Done():
			return c.JSON(http.StatusPartialContent, projects)
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
	id:=c.Param("taskID")
	taskID,err:=strconv.Atoi(id);if err!=nil{
		return c.JSON(http.StatusNotFound,"Invalid task ID")
	}
	t.TaskID=taskID

	v.Check(
		t.AssignedUsername.Valid,
		"user",
		"no valid user",
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
		if errors.Is(err,models.ErrRecordNotFound){
			c.Logger().Warn("Task not found :",err)
			return c.JSON(http.StatusNotFound,MapMessage("message",models.ErrRecordNotFound.Error()))
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
	r:=struct{
		Data 	[]*models.Pert 	`json:"data"`
		Result 	map[string]any	`json:"result"`
	}{}

	id:=c.Param("id")
	projectID,err:=strconv.Atoi(id);if err!=nil{
		return c.JSON(http.StatusNotFound,"Invalid project ID")
	}

	data,result,err:=calculate(&app.models.Pert,c.Request().Context(),projectID); if err!=nil{
		if errors.Is(err,models.ErrRecordNotFound){
			return c.JSON(http.StatusNotFound,MapMessage("message","No data CPM related data found"))
		}
		if errors.Is(err,ErrFetchingResult){
			r.Data=data
			return c.JSON(http.StatusPartialContent,r.Data)
		}
		c.Logger().Error("error getting cpm values : ",err)
		return c.JSON(http.StatusInternalServerError,MapMessage("error","Error getting CPM data"))
	}
	r.Data=data
	fmt.Println("r data",r.Data)
	r.Result=result.Result
	return c.JSON(http.StatusOK,r)
}

func(app *Application)CreatePert(c echo.Context)error{
	var pert []models.Pert	
	if err:=c.Bind(&pert);err!=nil{
		return c.JSON(http.StatusBadRequest,"Invalid request body")
	}
	if err:=app.models.Pert.Insert(c.Request().Context(),pert);err!=nil{
		c.Logger().Error(MapMessage("Pert Error",err.Error()))
		return c.JSON(http.StatusInternalServerError,MapMessage("error","failed to insert pert data"))
	}
	return c.JSON(http.StatusOK,map[string]string{"message":"pert data inserted successfully"})
}

func(app *Application)UpdatePert(c echo.Context)error{
	return c.JSON(http.StatusOK,"done")
}

func(app *Application)GetCpm(c echo.Context)error{
	r:=struct{
		Data 	[]*models.Cpm 	`json:"data"`
		Result 	map[string]any	`json:"result"`
	}{}

	id:=c.Param("id")
	projectID,err:=strconv.Atoi(id);if err!=nil{
		return c.JSON(http.StatusNotFound,"Invalid project ID")
	}

	data,result,err:=calculate(&app.models.Cpm,c.Request().Context(),projectID); if err!=nil{
		if errors.Is(err,models.ErrRecordNotFound){
			return c.JSON(http.StatusNotFound,MapMessage("message","No data CPM related data found"))
		}
		if errors.Is(err,ErrFetchingResult){
			r.Data=data
			return c.JSON(http.StatusPartialContent,r.Data)
		}
		c.Logger().Error("error getting cpm values : ",err)
		return c.JSON(http.StatusInternalServerError,MapMessage("error","Error getting CPM data"))
	}
	r.Data=data
	fmt.Println("r data",r.Data)
	r.Result=result.Result
	return c.JSON(http.StatusOK,r)
}

func(app *Application)CreateCpm(c echo.Context)error{
	var cpm []models.Cpm	
	if err:=c.Bind(&cpm);err!=nil{
		return c.JSON(http.StatusBadRequest,"Invalid request body")
	}
	if err:=app.models.Cpm.Insert(c.Request().Context(),cpm);err!=nil{
		c.Logger().Error(MapMessage("cpm Error",err.Error()))
		return c.JSON(http.StatusInternalServerError,MapMessage("error","failed to insert pert data"))
	}

	return c.JSON(http.StatusOK,MapMessage("message","cpm data inserted successfully"))
}

func(app *Application)UpdateCpm(c echo.Context)error{
	return c.JSON(http.StatusOK,"done")
}

func calculate[U models.Analytic,T models.ReadDatabase[U]](v T,ctx context.Context, id int)([]*U,models.Result,error){
	data,err:=v.GetData(ctx,id);if err!=nil{
		return nil,models.Result{},err 
	}
	if data==nil{
		return nil,models.Result{},models.ErrRecordNotFound
	}

	result,err:=v.GetResult(ctx,id); if err!=nil{
		if !errors.Is(err, models.ErrRecordNotFound){
			return data,models.Result{},ErrFetchingResult
		}
	}
	if result.Result!=nil{
		return data,models.Result{Result:result.Result},nil
	}
	result,err=external.RequestAndCalculatePERTCPM(data);if err!=nil{
		log.Println("Error fetching result: ",err)
		return data,models.Result{},ErrFetchingResult
	}

	log.Println("data ; ",data,"result ",result)
	return data,result,nil

}
