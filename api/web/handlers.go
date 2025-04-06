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
	"github.com/Raaffs/FluxMap/internal/sessionVar"
	validator "github.com/Raaffs/FluxMap/internal/validators"
	"github.com/gorilla/sessions"
	"github.com/guregu/null/v5"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo-contrib/session"
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


	sess,err:=session.Get(sessionvar.SESSION_NAME,c)
	if err!=nil{
		c.Logger().Error("error creating session: ",err)
		return c.JSON(http.StatusInternalServerError,"error creating session")
	}
	sess.Options=&sessions.Options{
			Path: "/",
			MaxAge: 86400*7,
			HttpOnly: true,
	}

	sess.Values[sessionvar.USERNAME]=u.Username
	log.Println("sess vals login",sess.Values)
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		c.Logger().Error("error saving session: ",err)
		return c.JSON(http.StatusInternalServerError,"error saving session")
	}

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
		validator.ErrPasswordTooWeak.Message,
	)

	v.Check(
		validator.Matches(u.Email,validator.EmailRX),
		"email",
		"invalid email",
	)

	if !v.Valid(){
		return c.JSON(http.StatusBadRequest,v)
	}

	hash,err:=HashPassword(u.Password);if err!=nil{
		return echo.NewHTTPError(echo.ErrInternalServerError.Code,"error creating user")
	}

	u.Created=time.Now().Format("2006-01-02")
	u.HashedPassword=hash
	if err:=app.models.Users.Create(context.Background(),u);err!=nil{
		if errors.Is(err,models.ErrAlreadyExist){
			return c.JSON(http.StatusConflict,MapMessage("error","User already exist"))
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

func (app *Application) CreateProject(c echo.Context) error {
	var p models.Project
	v := validator.New()

	// Bind JSON payload to the project struct
	if err := c.Bind(&p); err != nil {
		log.Println("Error json: ",err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON payload222",
		})
	}
	p.ProjectStartDate.NullTime.Time=time.Now()
	p.ProjectStartDate.NullTime.Valid=true
	// Validate project name length
	v.Check(
		validator.MinNameLength(p.ProjectName, 5),
		validator.ErrNameTooShort.Key,
		fmt.Sprintf(validator.ErrNameTooShort.Message, 5),
	)

	// Validate project description length
	v.Check(
		validator.MinDescriptionLength(p.ProjectDescription.String, 10),
		validator.ErrDescriptionTooShort.Key,
		fmt.Sprintf(validator.ErrDescriptionTooShort.Message, 10),
	)

	// Check if the username cookie exists
	sess, err := session.Get(sessionvar.SESSION_NAME,c);if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}
	username,ok:=sess.Values[sessionvar.USERNAME].(string);if !ok{
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	p.Ownername = username

	// If validation fails, return detailed validation errors
	if !v.Valid() {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation error",
			"fields": v.Errors, // Assuming v.Errors contains the validation error details
		})
	}

	// Attempt to create the project
	if err := app.models.Projects.Create(c.Request().Context(), p); err != nil {
		c.Logger().Error("Error creating project: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "An error occurred while creating the project",
		})
	}

	// Successfully created the project
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Project created successfully",
	})
}

func (app *Application) GetProjects(c echo.Context) error {
	sess, err := session.Get("session", c);if err != nil {
		return err
	}
	username,ok:=sess.Values[sessionvar.USERNAME].(string);if !ok{
		return c.JSON(http.StatusUnauthorized, map[string]string{"error":"you're not authorized"})
	}
	fmt.Println("username in get projects: ",username)
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

func(app *Application)GetAdminProjects(c echo.Context)error{
	sess, err := session.Get(sessionvar.SESSION_NAME, c);if err != nil {
		return err
	}
	fmt.Println("sess: ",sess.Values)

	username,ok:=sess.Values[sessionvar.USERNAME].(string);if !ok{
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	adminProjects,err:=app.models.Projects.RetrieveAdminProjects(c.Request().Context(),username); if err!=nil{
		c.Logger().Error("Error retrieving projects : ",err)
		c.JSON(http.StatusInternalServerError,MapMessage("Project","An error occurred while retrieving project"))
	}
	c.JSON(http.StatusOK,adminProjects)
	return nil
}

func(app *Application) GetManagerProjects(c echo.Context) error {
	sess, err := session.Get(sessionvar.SESSION_NAME,c);if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}
	username,ok:=sess.Values[sessionvar.USERNAME].(string);if !ok{
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	managerProjects, err := app.models.Projects.RetrieveManagerProjects(c.Request().Context(), username)
	if err != nil {
		c.Logger().Error("Error retrieving manager projects: ", err)
		return c.JSON(http.StatusInternalServerError, MapMessage("Project", "An error occurred while retrieving manager projects"))
	}	

	return c.JSON(http.StatusOK, managerProjects)
}

func(app *Application) GetAssignedProjects(c echo.Context) error {
	user, err := c.Cookie("username")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	assignedProjects, err := app.models.Projects.RetrieveAssginedProjects(c.Request().Context(), user.Value)
	if err != nil {
		// Log and handle errors
		c.Logger().Error("Error retrieving assigned projects: ", err)
		return c.JSON(http.StatusInternalServerError, MapMessage("Project", "An error occurred while retrieving assigned projects"))
	}
	return c.JSON(http.StatusOK, assignedProjects)
}

func(app *Application)GetProjectByID(c echo.Context)error{
	id, err := strconv.Atoi(c.Param("id"));if err!=nil{
		return c.JSON(http.StatusBadRequest,MapMessage("message","Invalid project id"))
	}
	projects,err:=app.models.Projects.RetrieveProjectByID(c.Request().Context(),id);if err!=nil{
		if errors.Is(err,models.ErrRecordNotFound){
			return c.JSON(http.StatusNotFound,MapMessage("message","Project not found"))
		}
		c.Logger().Error("Error retrieving project by id: ",err)
		return c.JSON(http.StatusInternalServerError,MapMessage("message","An error occurred while retrieving project by id"))
	}
	return c.JSON(http.StatusOK,projects)
}


func(app *Application)Invite(c echo.Context)error{
	invitation:=struct{
		Username  string `json:"username"`
		Role	  string `json:"role"`
	}{}

	id, err := strconv.Atoi(c.Param("id"));if err!=nil{
		return c.JSON(http.StatusBadRequest,map[string]string{"error":"invalid project id"})
	}

	if err:=c.Bind(&invitation);err!=nil{
		return c.JSON(http.StatusBadRequest,map[string]string{"error":"invalid json"})
	}

	fmt.Println("username, role: ",invitation.Username,invitation.Role)
	exist,err:=app.models.Users.Exist(c.Request().Context(),invitation.Username)
	if !exist{
		c.Logger().Error("user doesn't exists")
		return c.JSON(http.StatusNotFound,map[string]string{"error":"user doesn't exist"})
	}

	if err!=nil{
		c.Logger().Error("Error inviting user:",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
	}

	alreadyInvited,err:=app.models.Invitation.Exist(c.Request().Context(),invitation.Username,id); if err!=nil{
		c.Logger().Error("Error checking invitation:",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
	}
	if alreadyInvited{
		return c.JSON(http.StatusConflict,map[string]string{"error":"user is already invited"})
	}

	if err:=app.models.Invitation.Invite(c.Request().Context(),invitation.Username,id,invitation.Role); err!=nil{
		c.Logger().Error("Error inviting user: ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
	}
	return nil
}



func (app *Application)GetInvitations(c echo.Context)error{
	sess,err:=session.Get("session",c);if err!=nil{
		c.Logger().Error("error getting session: ",err)
		return c.JSON(http.StatusUnauthorized,MapMessage("session","Error getting session"))
	}

	username,ok:=sess.Values[sessionvar.USERNAME].(string)
	if !ok{
		c.Logger().Error("Error getting username from session")
		return c.JSON(http.StatusUnauthorized,MapMessage("session","Error getting username from session"))
	}

	invitations,err:=app.models.Invitation.GetInvitations(c.Request().Context(),username);if err!=nil{
		
		c.Logger().Error("Error getting invitations: ",err)
		return c.JSON(http.StatusInternalServerError,MapMessage("Invitations","Failed to get invitations"))
	}
	return c.JSON(http.StatusOK,invitations)
}

func (app *Application) ConfirmInvitation(c echo.Context) error {

	invitation:=struct{
		Status string `json:"status"`
	}{}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error("Error converting id to int: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}
	
	if err:=c.Bind(&invitation);err!=nil{
		c.Logger().Error("Error binding request body: ",err)
		return c.JSON(http.StatusBadRequest,map[string]string{"error":"Invalid request body"})
	}

	if invitation.Status!="rejected" && invitation.Status!="accepted"{
		return c.JSON(http.StatusBadRequest,map[string]string{"error":"Invalid invitation status"})
	}

	sess, err := session.Get("session", c)
	if err != nil {
		c.Logger().Error("Error getting session: ", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Error getting session"})
	}

	username, ok := sess.Values[sessionvar.USERNAME].(string)
	if !ok {
		c.Logger().Error("Error retrieving username from session")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Error getting username from session"})
	}

	inviteID, err := app.models.Invitation.ConfirmInvitation(c.Request().Context(), invitation.Status, id, username)
	if err != nil {
		c.Logger().Error("Error confirming invitation: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to confirm invitation"})
	}
	log.Println("invite returningid: ",inviteID)
	invitationDetail, err := app.models.Invitation.GetInvitationByID(c.Request().Context(), inviteID)
	log.Println("invitation details: ",invitationDetail)
	if err != nil {
		c.Logger().Error("Error retrieving invitation detail: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get invitation detail"})
	}

	if invitationDetail.Role == "manager" {
		log.Println("added to manager")
		if err := app.models.Projects.AssignManager(c.Request().Context(), username, invitationDetail.ProjectID); err != nil {
			c.Logger().Error("Error assigning manager: ", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign manager role"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Invitation confirmed successfully"})
}


func (app *Application)CreateTask(c echo.Context)error{
	var t models.Task
	v:= validator.New()
	if err:=c.Bind(&t);err!=nil{
		return c.JSON(http.StatusBadRequest,map[string]string{"error":"Invalid Json payload"})
	}
	id,err:=strconv.Atoi(c.Param("id"));if err!=nil{
		c.Logger().Error(MapMessage("error converting to string",err.Error()))
		return c.JSON(http.StatusBadRequest,map[string]string{"error":"Invalid project ID"})
	}
	t.ParentProjectID=id
	v.Check(
		validator.MinNameLength(t.TaskName,3),
		validator.ErrNameTooShort.Key,
		validator.ErrNameTooShort.Message,
	)
	if !v.Valid(){
		return c.JSON(http.StatusBadRequest,map[string]any{"error":v.Errors})
	}
	sess, err := session.Get("session", c);if err != nil {
		return err
	}
	username,ok:=sess.Values[sessionvar.USERNAME].(string);if !ok{
		return c.JSON(http.StatusUnauthorized, map[string]string{"error":"you're not authorized"})
	}
	t.Createdby=username
	if err:=app.models.Task.Create(c.Request().Context(),t);err!=nil{
		c.Logger().Error(MapMessage("Error creating task",err.Error()))
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"An error occured"})
	}
	return c.JSON(http.StatusOK,"task created")
}

func (app *Application)GetTasks(c echo.Context)error{
	taskfunc:= app.GetTaskBasedOnAccess(app.ManagerTasks,app.UserTasks)
	return taskfunc(c)
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

func(app *Application)UpdateUserTask(c echo.Context)error{
	status:=struct{
		TaskStatus string `json:"taskStatus"`
	}{}
	v:=validator.New()

	id,err:=strconv.Atoi(c.Param("taskID")); if err!=nil{
		return c.JSON(http.StatusNotFound,MapMessage("error","task not found"))
	}
	log.Println("invalid task id ",id)
	if err:=c.Bind(&status);err!=nil{
		return c.JSON(http.StatusBadRequest,MapMessage("error","invalid request"))
	}
	v.Check(
		status.TaskStatus=="completed"||status.TaskStatus=="pending",
		"status",
		"invalid status",
	)
	log.Println("invalid status")

	if err:=app.models.Task.UpdateTask(c.Request().Context(),id,status.TaskStatus);err!=nil{
		c.Logger().Error("error updating task: ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
	}

	return nil
}

func (app *Application)UpdateManagerTask(c echo.Context)error{
	var t models.Task
	v:=validator.New()
	if err:=c.Bind(&t);err!=nil{
		return c.JSON(http.StatusBadRequest,MapMessage("error","invalid request"))
	}
	v.Check(
		validator.MinNameLength(t.TaskName,3),
		validator.ErrNameTooShort.Key,
		validator.ErrNameTooShort.Message,
	)
	if !v.Valid(){
		return c.JSON(http.StatusBadRequest,v)
	}
	if !t.AssignedUsername.Valid{
		return c.JSON(http.StatusBadRequest,"task must be assigned to a user")
	}
	
	if t.Approved.Bool{
		t.TaskApprovedDate.Time =time.Now()
	}
	if err:=app.models.Task.UpdateManagerTask(c.Request().Context(),t);err!=nil{
		c.Logger().Error(MapMessage("Error updating task",err.Error()))
		return c.JSON(http.StatusInternalServerError,MapMessage("error","internal server error"))
	}
	return c.JSON(http.StatusOK,"updated successfully")
}

func(app *Application)AddManager(c echo.Context)error{
	m:=struct{
		Manager string `json:"manager"`
	}{}
	projectID,err:=strconv.Atoi(c.Param("id"));if err!=nil{
		c.Logger().Error(MapMessage("error converting to string",err.Error()))
		return c.JSON(http.StatusBadRequest,map[string]string{"error":"Invalid project ID"})
	}
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
	log.Println("t.assigned user: ",t.AssignedUsername,t.AssignedUsername.Valid)

	v:=validator.New()
	if err:=c.Bind(&t);err!=nil{
		c.Logger().Error("error reading json,",err)
		return c.JSON(http.StatusBadRequest,"Invalid request body")
	}
	log.Println("t.assigned user: ",t.AssignedUsername,t.AssignedUsername.Valid)
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

	if t.Approved.ValueOrZero(){
		t.TaskApprovedDate=null.NewTime(time.Now(),true)
	}else{
		t.TaskApprovedDate=null.NewTime(time.Time{},false)
	}
	
	if !v.Valid(){
		c.Logger().Error(v)
		return c.JSON(http.StatusBadRequest,v)
	}
	log.Println("tasks: ",t)
	if err:=app.models.Task.UpdateManagerTask(c.Request().Context(),t);err!=nil{
		if errors.Is(err,models.ErrRecordNotFound){
			c.Logger().Warn("Task not found :",err)
			return c.JSON(http.StatusNotFound,MapMessage("message",models.ErrRecordNotFound.Error()))
		}
		c.Logger().Error("error updating manager task : ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"failed to update task"})
	}
	log.Println("before deferred")
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

	data,result,err:=GetAnalytics(&app.models.Pert,c.Request().Context(),projectID); if err!=nil{
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
		c.Logger().Error("error binding pert : ",err)
		return c.JSON(http.StatusBadRequest,map[string]string{"error":"invalid json body"})
	}
	id:=c.Param("id")
	log.Println("HTE OF PROJE ",id)
	projectID,err:=strconv.Atoi(id);if err!=nil{
		return c.JSON(http.StatusNotFound,"Invalid project ID")
	}
	pert[0].ParentProjectID=projectID
	log.Println("new pert inseret: ",pert)

	if err:=app.models.Pert.Insert(c.Request().Context(),pert);err!=nil{
		c.Logger().Error(MapMessage("Pert Error",err.Error()))
		return c.JSON(http.StatusInternalServerError,MapMessage("error","failed to insert pert data"))
	}
	log.Println("new pert inseret: ",pert)
	if err:=Calculate[models.Pert,*models.PertModel[models.Pert]](&app.models.Pert,c.Request().Context(),pert[0].ParentProjectID);err!=nil{
		c.Logger().Error("Error calculating pert : ",err)
		return c.JSON(http.StatusInternalServerError,MapMessage("Error","Failed to calculate values"))
	}
	return c.JSON(http.StatusOK,MapMessage("PERT","data and result inserted successfully"))
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

	data,result,err:=GetAnalytics(&app.models.Cpm,c.Request().Context(),projectID); if err!=nil{
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
	r.Result=result.Result
	return c.JSON(http.StatusOK,r)
}

func(app *Application)CreateCpm(c echo.Context)error{
	var cpm []models.Cpm	
	if err:=c.Bind(&cpm);err!=nil{
		return c.JSON(http.StatusBadRequest,"Invalid request body")
	}
	if err:=DetectCycleCpm(cpm);err!=nil{
		c.Logger().Warn(err)
		return c.JSON(http.StatusBadRequest,"Cyclic dependencies are not allowed")
	}
	if err:=app.models.Cpm.Insert(c.Request().Context(),cpm);err!=nil{
		c.Logger().Error(MapMessage("cpm Error",err.Error()))
		return c.JSON(http.StatusInternalServerError,MapMessage("error","failed to insert cpm data"))
	}
	fmt.Println("after insert")

	if err:=Calculate(&app.models.Cpm,c.Request().Context(),cpm[0].ParentProjectID); err!=nil{
		c.Logger().Error("Error calculating cpm data : ",err)
		return c.JSON(http.StatusInternalServerError,MapMessage("error","Failed to calculate CPM values"))
	}
	return c.JSON(http.StatusOK,MapMessage("message","cpm data inserted successfully"))
}

func(app *Application)UpdateCpm(c echo.Context)error{
	return c.JSON(http.StatusOK,"done")
}




