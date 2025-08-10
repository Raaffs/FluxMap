package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Raaffs/FluxMap/internal/models"
	sessionvar "github.com/Raaffs/FluxMap/internal/sessionVar"
	validator "github.com/Raaffs/FluxMap/internal/validators"
	"github.com/gorilla/sessions"
	"github.com/guregu/null/v5"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (app *Application) Login(c echo.Context) error {
	var u models.User
	err := c.Bind(&u)
	if err != nil {
		c.Logger().Error("error binding json : ", err)
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "Invalid credentials"})
	}
	u.Username = strings.TrimSpace(u.Username)
	if err := app.models.Users.Login(c.Request().Context(), u.Username, u.Password); err != nil {
		if errors.Is(err, models.ErrInvalidCredential) {
			c.Logger().Warn("invalid auth")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		c.Logger().Error("Error authenticating user: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to log in"})
	}

	sess, err := session.Get(sessionvar.SESSION_NAME, c)
	if err != nil {
		c.Logger().Error("error creating session: ", err)
		return c.JSON(http.StatusInternalServerError, "error creating session")
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	sess.Values[sessionvar.USERNAME] = u.Username
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		c.Logger().Error("error saving session: ", err)
		return c.JSON(http.StatusInternalServerError, "error saving session")
	}
	roleToProjectMap,err:= app.CacheUserProjectsToSession(c); if err!=nil{
		c.Logger().Error("Error caching projects : ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"interval server error"})
	}
	app.SendUpdateNotification(c.Request().Context(),u.Username)
	return c.JSON(http.StatusOK, map[string]any{"roles":roleToProjectMap})
}

func (app *Application)SendProjectRoleCache(c echo.Context)error{
	roleToProjectMap,err:=app.CacheUserProjectsToSession(c);if err!=nil{
		c.Logger().Error("Error caching user projects: ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
	}
	return c.JSON(http.StatusOK,map[string]any{"roles":roleToProjectMap})
}

func (app *Application) Register(c echo.Context) error {
	u := models.User{}
	v := validator.New()
	if err := c.Bind(&u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	v.Check(
		validator.MinNameLength(u.Username),
		validator.ErrNameTooShort.Key,
		validator.ErrNameTooShort.Message,
	)

	v.Check(
		validator.IsStrongPassword(u.Password),
		validator.ErrPasswordTooWeak.Key,
		validator.ErrPasswordTooWeak.Message,
	)

	v.Check(
		validator.Matches(u.Email, validator.EmailRX),
		"email",
		"invalid email",
	)

	if !v.Valid() {
		return c.JSON(http.StatusBadRequest, v)
	}

	hash, err := HashPassword(u.Password)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "error creating user")
	}

	u.Created = time.Now().Format("2006-01-02")
	u.HashedPassword = hash
	if err := app.models.Users.Create(context.Background(), u); err != nil {
		if errors.Is(err, models.ErrAlreadyExist) {
			return c.JSON(http.StatusConflict, map[string]string{"error":"user already exist"})
		}
		c.Logger().Error("Error creating user: ", err)
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "internal server error")
	}

	return c.JSON(http.StatusOK, "user registered successfully")
}



func (app *Application) SessionCheck(c echo.Context) error {
	sess, err := session.Get(sessionvar.SESSION_NAME, c)
	if err != nil {
		c.Logger().Error("failed to get session: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]bool{"isAuthenticated": false})
	}

	username, ok := sess.Values[sessionvar.USERNAME].(string)
	if !ok || username == "" {
		return c.JSON(http.StatusOK, map[string]bool{"isAuthenticated": false})
	}

	return c.JSON(http.StatusOK, map[string]bool{
		"isAuthenticated": true,
	})
}

func (app *Application) Logout(c echo.Context) error {
	session, err := session.Get("session", c)
	if err != nil {
		c.Logger().Error("Error logging out : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	session.Options.MaxAge = -1
	if err = session.Save(c.Request(), c.Response()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, "")
}

func (app *Application) CreateProject(c echo.Context) error {
	var p models.Project
	v := validator.New()

	// Bind JSON payload to the project struct
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON payload",
		})
	}

	p.ProjectStartDate.NullTime.Time = time.Now()
	p.ProjectStartDate.NullTime.Valid = true
	// Validate project name length
	v.Check(
		validator.MinNameLength(p.ProjectName),
		validator.ErrNameTooShort.Key,
		fmt.Sprintf(validator.ErrNameTooShort.Message, 5),
	)

	// Validate project description length
	v.Check(
		validator.MinDescriptionLength(p.ProjectDescription.String),
		validator.ErrDescriptionTooShort.Key,
		fmt.Sprintf(validator.ErrDescriptionTooShort.Message, 10),
	)

	// Check if the username session exists
	username := c.Get(sessionvar.USERNAME).(string)
	p.Ownername = username

	// If validation fails, return detailed validation errors
	if !v.Valid() {
		log.Println("project errors: ",v.Errors)
		return c.JSON(http.StatusBadRequest, v)
	}

	// Attempt to create the project
	if err := app.models.Projects.Create(c.Request().Context(), p); err != nil {
		c.Logger().Error("Error creating project: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "An error occurred while creating the project",
		})
	}
	roleMap,err:=app.CacheUserProjectsToSession(c);if err!=nil{
		c.Logger().Error("Error caching projects to session: ",err)
	}
	return c.JSON(http.StatusOK, map[string]any{"roles":roleMap})
}

func (app *Application) GetProjects(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	username, ok := sess.Values[sessionvar.USERNAME].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "you're not authorized"})
	}

	resultChan := make(chan ProjectResult, 1)
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	go app.FetchProjects(ctx, username, resultChan)
	result := <-resultChan
	
	if result.Err != nil {
		if errors.Is(result.Err, context.DeadlineExceeded) {
			return c.JSON(http.StatusPartialContent, result)
		}
		c.Logger().Error("Error fetching projects: ", result.Err)
		return c.JSON(http.StatusInternalServerError, "error retrieving projects")
	}
	return c.JSON(http.StatusOK, result)
}

func (app *Application) GetAdminProjects(c echo.Context) error {
	username := c.Get(sessionvar.USERNAME).(string)

	adminProjects, err := app.models.Projects.RetrieveAdminProjects(c.Request().Context(), username)
	if err != nil {
		c.Logger().Error("Error retrieving projects : ", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error":"internal server error"})
	}
	c.JSON(http.StatusOK, adminProjects)
	return nil
}

func (app *Application) GetManagerProjects(c echo.Context) error {
	username := c.Get(sessionvar.USERNAME).(string)

	managerProjects, err := app.models.Projects.RetrieveManagerProjects(c.Request().Context(), username)
	if err != nil {
		c.Logger().Error("Error retrieving manager projects: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error":"an error occurred while retrieving projects"})
	}

	return c.JSON(http.StatusOK, managerProjects)
}

func (app *Application) GetAssignedProjects(c echo.Context) error {
	username := c.Get(sessionvar.USERNAME).(string)

	assignedProjects, err := app.models.Projects.RetrieveAssginedProjects(c.Request().Context(), username)
	if err != nil {
		c.Logger().Error("Error retrieving assigned projects: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error":"An error occurred while retrieving assigned projects"})
	}
	return c.JSON(http.StatusOK, assignedProjects)
}

func (app *Application) GetProjectByID(c echo.Context) error {
	id := c.Param("id")
	projID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid project id"})
	}

	projects, err := app.models.Projects.RetrieveProjectByID(c.Request().Context(), projID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Project not found"})
		}
		c.Logger().Error("Error retrieving project by id: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An error occurred while retrieving project by id"})
	}
	return c.JSON(http.StatusOK, projects)
}

func (app *Application) Invite(c echo.Context) error {
	invitation := struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}{}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project id"})
	}

	if err := c.Bind(&invitation); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
	}

	switch status := app.EnsureExists(c, func(ctx context.Context) (bool, error) {
		return app.models.Users.Exist(ctx, invitation.Username)
	}); status {
	case ErrorCheckingExistStatus:
	// error response already sent by EnsureExists, just exit
		return nil	
	case NotExists:
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user doesn't exist"})
	case Exists:
		//excute rest of the code if user exists 
	}

	switch status:=app.EnsureExists(c,func(ctx context.Context) (bool, error) {
		return app.models.Invitation.Exist(ctx,invitation.Username,id)
	});status{
	case ErrorCheckingExistStatus:
	// error response already sent by EnsureExists, just exit
		return nil
	case Exists: 
		return c.JSON(http.StatusConflict,map[string]string{"error":"user already invited"})
	case NotExists:
		//continue with inviting user
	}

	if err := app.models.Invitation.Invite(c.Request().Context(), invitation.Username, id, invitation.Role); err != nil {
		c.Logger().Error("Error inviting user: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	app.SendInviteNotification(c.Request().Context(),invitation.Username)

	return c.JSON(http.StatusOK, map[string]string{"message": "invite sent"})
}

func (app *Application) GetInvitations(c echo.Context) error {
	username := c.Get(sessionvar.USERNAME).(string)
	invitations, err := app.models.Invitation.GetPendingInvitations(c.Request().Context(), username)
	if err != nil {

		c.Logger().Error("Error getting invitations: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"Invitations": "Failed to get invitations"})
	}
	if err:=app.models.Invitation.SetRead(c.Request().Context(),username);err!=nil{
		c.Logger().Error("Error updating status of hasread column: ",err)
	}
	return c.JSON(http.StatusOK, invitations)
}

func (app *Application) ConfirmInvitation(c echo.Context) error {
	username := c.Get(sessionvar.USERNAME).(string)

	invitation := struct {
		Status string `json:"status"`
	}{}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Warn("Error converting id to int: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	if err := c.Bind(&invitation); err != nil {
		c.Logger().Warn("Error binding request body: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if invitation.Status != "rejected" && invitation.Status != "accepted" {
		c.Logger().Warn("invalid status recieved : ", invitation.Status)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid invitation status"})
	}

	inviteID, err := app.models.Invitation.ConfirmInvitation(c.Request().Context(), invitation.Status, id, username)
	if err != nil {
		c.Logger().Error("Error confirming invitation: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to confirm invitation"})
	}

	invitationDetail, err := app.models.Invitation.GetInvitationByID(c.Request().Context(), inviteID)
	if err != nil {
		c.Logger().Error("Error retrieving invitation detail: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get invitation detail"})
	}

	if invitationDetail.Role == "manager" {
		if err := app.models.Projects.AssignManager(c.Request().Context(), username, invitationDetail.ProjectID); err != nil {
			c.Logger().Error("Error assigning manager: ", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign manager role"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Invitation confirmed successfully"})
}

func(app *Application)RemoveUser(c echo.Context)error{
	removeUser:=struct{
		Username string `json:"username"`
		FallBackUser string `json:"fallBackUser"`
	}{}
	projectID,err:=strconv.Atoi(c.Param("id"));if err!=nil{
		c.Logger().Warn("Error converting projectid string to int: ",err)
		return c.JSON(http.StatusNotFound,map[string]string{"error":"project not found"})

	}
	if err:=c.Bind(&removeUser);err!=nil{
		c.Logger().Error("Error binding json to removeUser struct: ",err)
		return c.JSON(http.StatusBadRequest,map[string]string{"error":"invalid json format"})
	}

	if err:=app.models.Task.ReallocateUser(c.Request().Context(),projectID,removeUser.Username,removeUser.FallBackUser);err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return c.JSON(http.StatusNotFound,map[string]string{"error":"user not found"})
		}
		c.Logger().Error("Error transferring task to user: ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
	}

	if err:=app.models.Invitation.Delete(c.Request().Context(),projectID,removeUser.Username);err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return c.JSON(http.StatusNotFound,map[string]string{"error":"user not found"})
		}
		c.Logger().Error("Error removing user from invitation: ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
	}
	return c.JSON(http.StatusOK,map[string]string{"message":"user removed successfully"})
}

func (app *Application)CreateTask(c echo.Context) error {
	var t models.Task
	username := c.Get(sessionvar.USERNAME).(string)
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Json payload"})
	}

	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Warn("error converting to string", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}
	if accepted:=app.CheckInvitationStatus(c,t,projectID,username);!accepted{
		return nil
	}

	v:=ValidateTask(t)
	if !v.Valid(){
		c.Logger().Error("invalid task : ",v.Errors)
		return c.JSON(http.StatusBadRequest,v)	
	}

	t.ParentProjectID = projectID
	t.Createdby = username
	if err := app.models.Task.Create(c.Request().Context(), t); err != nil {
		c.Logger().Error("Error creating task", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "An error occured"})
	}
	msg:=fmt.Sprintf("Task %s has been assigned to you by %s",t.TaskName,t.Createdby)
	if err:=app.UpdateCountAndNotify(c,projectID,msg,t.AssignedUsername.String);err!=nil{
		c.Logger().Error("Error sending notification: ",err)
	}
	return c.JSON(http.StatusOK, "task created")
}

func (app *Application) GetTasks(c echo.Context) error {
	taskfunc := app.GetTaskBasedOnAccess(app.ManagerTasks, app.UserTasks)
	return taskfunc(c)
}

func (app *Application) GetTaskByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.Logger().Error("error converting to string", err.Error())
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}
	task, err := app.models.Task.GetByID(c.Request().Context(), id)
	if err != nil {
		c.Logger().Error("Error retrieving task", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve task"})
	}

	return c.JSON(http.StatusOK, task)
}

func (app *Application)GetOverDueTasks(c echo.Context)error{
	username:=c.Get(sessionvar.USERNAME).(string)
	tasks,err:=app.models.Task.GetOverdue(c.Request().Context(),username);if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return c.JSON(http.StatusNotFound,map[string]string{"message":"no overdue tasks found"})
		}
		c.Logger().Error("Error retrieving overdue tasks: ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
	}
	log.Println("tasks: ",tasks)
	return c.JSON(http.StatusOK,map[string]any{"overdue":tasks})
}

func (app *Application) GetConfirmedUsers(c echo.Context) error {
	username := c.Get(sessionvar.USERNAME).(string)
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Warn("Project with id %d not found ", projectID, " err: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid projectID"})
	}
	users, err := app.models.Invitation.FetchConfirmedMembers(c.Request().Context(), projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Logger().Warn("No users found for project with id %d ", projectID)
			return c.JSON(http.StatusPartialContent, map[string]string{"users": username})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "An error occurred while retrieving projects"})
	}
	//appened yourself to project if you're the creator of project
	users = append(users, &username)
	return c.JSON(http.StatusOK, map[string][]*string{"users": users})
}

func (app *Application) GetAssignedUserByTask(c echo.Context) error {
	return c.JSON(http.StatusOK, "assgined user by task retrived")
}

func (app *Application) UpdateProject(c echo.Context) error {
	return c.JSON(http.StatusOK, "updated successfully")
}

func (app *Application) UpdateUserTask(c echo.Context) error {
	status := struct {
		TaskStatus string `json:"taskStatus"`
	}{}
	username:=c.Get(sessionvar.USERNAME).(string)
	v := validator.New()
	taskID, err := strconv.Atoi(c.Param("taskID"));if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}
	projectID,err:=strconv.Atoi(c.Param("id"));if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}
	if err := c.Bind(&status); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	v.Check(
		status.TaskStatus == "completed" || status.TaskStatus == "pending",
		"status",
		"invalid status",
	)
	if !v.Valid(){
		c.Logger().Warn("invalid status : ",status)
		return c.JSON(http.StatusBadRequest,map[string]string{"error":"invalid status"})
	}
	assignedByName, err := app.models.Task.UpdateStatus(c.Request().Context(), projectID, taskID, status.TaskStatus);if err != nil {
		c.Logger().Error("error updating task: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	msg:=fmt.Sprintf("Task status set to %s by %s",status.TaskStatus,username)
	if err:=app.UpdateCountAndNotify(c,projectID,msg,assignedByName);err!=nil{
		c.Logger().Error("Error sending update notification: ",err)
	}

	return nil
}

func (app *Application)ApproveTask(c echo.Context) error {
	var t models.Task

	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Warn("Error converting project id from string to int : ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	taskID, err := strconv.Atoi( c.Param("taskID"))
	if err != nil {
		return c.JSON(http.StatusNotFound, "Invalid task ID")
	}
	if err := c.Bind(&t); err != nil {
		c.Logger().Warn("Error reading json: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	t.TaskID = taskID

	if accepted:=app.CheckInvitationStatus(c,t,projectID,t.AssignedUsername.String);!accepted{
		return nil
	}

	v := ValidateManagerUpdate(t)
	if !v.Valid() {
		c.Logger().Error("Validation failed for task: ", v.Errors)
		return c.JSON(http.StatusBadRequest, v)
	}

	if t.Approved.ValueOrZero() {
		t.TaskApprovedDate = null.NewTime(time.Now(), true)
	} else {
		t.TaskApprovedDate = null.NewTime(time.Time{}, false)
	}

	msg,err:=app.GenerateUpdateMessage(c.Request().Context(),t);if err!=nil{
		c.Logger().Error("error generating update message : ", err)
	}

	if err := app.models.Task.ManagerAuthorizedUpdate(c.Request().Context(),projectID, t); err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			c.Logger().Warn("Task not found :", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
		}
		c.Logger().Error("error approving task : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update task"})
	}

	if err:=app.UpdateCountAndNotify(c,projectID,msg,t.AssignedUsername.String);err!=nil{
		c.Logger().Error("Error sending update notification: ",err)
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "task approved successfully"})
}

func (app *Application)ManagerRestrictedTask(c echo.Context) error {
	var t models.Task
	projectID,err:=strconv.Atoi(c.Param("id"));if err!=nil{
		return c.JSON(http.StatusNotFound, "Invalid task ID")
	}

	if err := c.Bind(&t); err != nil {
		c.Logger().Error("error reading json,", err)
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	id := c.Param("taskID");
	taskID, err := strconv.Atoi(id); if err != nil {
		return c.JSON(http.StatusNotFound, "Invalid task ID")
	}

	t.TaskID = taskID
	v:=ValidateManagerUpdate(t)
	if !v.Valid() {
		c.Logger().Error(v)
		return c.JSON(http.StatusBadRequest, v.Errors)
	}

	if accepted:=app.CheckInvitationStatus(c,t,projectID,t.AssignedUsername.String);!accepted{
		return nil
	}

	msg,err:=app.GenerateUpdateMessage(c.Request().Context(),t);if err!=nil{
		c.Logger().Error("error generating update message : ", err)
	}

	if err := app.models.Task.ManagerAuthorizedUpdate(c.Request().Context(),projectID, t); err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			c.Logger().Warn("Task not found :", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error":"task not found"})
		}
		c.Logger().Error("error updating manager task : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update task"})
	}
	
	if err:=app.UpdateCountAndNotify(c,projectID,msg,t.AssignedUsername.String);err!=nil{
		c.Logger().Error("Error sending update notification: ",err)
	}

	return c.JSON(http.StatusOK, "task updated")
}

func (app *Application) GetPert(c echo.Context) error {
	r := struct {
		Data   []*models.Pert `json:"data"`
		Result map[string]any `json:"result"`
	}{}

	projectID, err := strconv.Atoi(c.Param("id")); if err != nil {
		return c.JSON(http.StatusNotFound, "Invalid project ID")
	}

	data, result, err := GetAnalytics(&app.models.Pert, c.Request().Context(), projectID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "No data CPM related data found"})
		}
		if errors.Is(err, ErrFetchingResult) {
			r.Data = data
			return c.JSON(http.StatusPartialContent, r.Data)
		}
		c.Logger().Error("error getting cpm values : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error getting CPM data"})
	}
	r.Data = data
	r.Result = result.Result

	return c.JSON(http.StatusOK, r)
}

func (app *Application) CreatePert(c echo.Context) error {
	var pert []models.Pert
	if err := c.Bind(&pert); err != nil {
		c.Logger().Error("error binding pert : ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json body"})
	}
	
	id := c.Param("id")
	projectID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Invalid project ID")
	}
	if len(pert) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no pert data provided"})
	}

	pert[0].ParentProjectID = projectID

	if err := app.models.Pert.Insert(c.Request().Context(), pert); err != nil {
		c.Logger().Error("error inserting pert data : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to insert pert data"})
	}
	if err := Calculate[models.Pert, *models.PertModel[models.Pert]](&app.models.Pert, c.Request().Context(), pert[0].ParentProjectID); err != nil {
		c.Logger().Error("Error calculating pert : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"Error": "Failed to calculate values"})
	}
	return c.JSON(http.StatusOK, map[string]string{"PERT": "data and result inserted successfully"})
}

func (app *Application) GetCpm(c echo.Context) error {
	r := struct {
		Data   []*models.Cpm  `json:"data"`
		Result map[string]any `json:"result"`
	}{}

	id := c.Param("id")
	projectID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Invalid project ID")
	}

	data, result, err := GetAnalytics(&app.models.Cpm, c.Request().Context(), projectID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "No data CPM related data found"})
		}
		if errors.Is(err, ErrFetchingResult) {
			r.Data = data
			return c.JSON(http.StatusPartialContent, r.Data)
		}
		c.Logger().Error("error getting cpm values : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error getting CPM data"})
	}
	r.Data = data
	r.Result = result.Result
	return c.JSON(http.StatusOK, r)
}

func (app *Application) CreateCpm(c echo.Context) error {
	var cpm []models.Cpm
	if err := c.Bind(&cpm); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if len(cpm) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no cpm data provided"})
	}
	log.Println("cpm recieved data: ",cpm[0])
	if err := DetectCycleCpm(cpm); err != nil {
		c.Logger().Warn(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cyclic dependencies are not allowed"})
	}
	if err := app.models.Cpm.Insert(c.Request().Context(), cpm); err != nil {
		c.Logger().Error(map[string]string{"error": "failed to insert cpm data"})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to insert cpm data"})
	}
	if err := Calculate(&app.models.Cpm, c.Request().Context(), cpm[0].ParentProjectID); err != nil {
		c.Logger().Error(map[string]string{"error": "Failed to calculate CPM values"})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to calculate CPM values"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "cpm data inserted successfully"})
}

func (app *Application) GetAllUpdates(c echo.Context) error {
	username := c.Get(sessionvar.USERNAME).(string)
	updates, err := app.models.Update.GetAll(c.Request().Context(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Logger().Warn("not updates : ", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No updates found"})
		}
		c.Logger().Error("error getting updates : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	if err:=app.models.Update.SetRead(c.Request().Context(),username);err!=nil{
		c.Logger().Error("Error updating status of hasread column: ",err)
	}
	count,err:=app.models.Update.GetTotalUnread(c.Request().Context(),username);if err!=nil{
		c.Logger().Error("Error getting total updates: ",err)
	}

	errchan:=app.websocket.PushToClients([]byte(strconv.Itoa(count)),
		func(w WSUser) bool {
		return w.Username==username
	})

	app.CheckChannelError(errchan)
	return c.JSON(http.StatusOK, updates)
}

func (app *Application) GetRecentUpdates(c echo.Context) error {
	username := c.Get(sessionvar.USERNAME).(string)
	updates, err := app.models.Update.GetRecent(c.Request().Context(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Logger().Warn("not updates : ", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No updates found"})
		}
		c.Logger().Error("error getting updates : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, updates)
}


func (app *Application)RemoveTask(c echo.Context)error{
	taskID,err:=strconv.Atoi(c.Param("taskID"));if err!=nil{
		return c.JSON(http.StatusBadRequest,map[string]string{"error":"Invalid task id"})
	}
	projectID,err:=strconv.Atoi(c.Param("id"));if err!=nil{
		return c.JSON(http.StatusBadRequest,map[string]string{"error":"Invalid task id"})
	}
	if err:=app.models.Task.Archieve(c.Request().Context(),projectID,taskID); err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			c.Logger().Errorf("Task with %d id not found. Error: %s",taskID,err)
			return c.JSON(http.StatusNotFound,map[string]string{"error":"task not found"})
		}
		c.Logger().Error("Error deleting task : ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
	}
	return c.JSON(http.StatusOK,map[string]string{"message":"Task successfully removed"})
}


func (app *Application) GetProjectUpdates(c echo.Context) error {
	username := c.Get(sessionvar.USERNAME).(string)
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Warn("invalid project id")
		return c.JSON(http.StatusBadRequest, "Invalid project id")
	}
	updates, err := app.models.Update.GetProjectUpdates(c.Request().Context(), projectID, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Logger().Warn("not updates : ", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No updates found"})
		}
		c.Logger().Error("error getting updates : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	if err:=app.models.Update.SetRead(c.Request().Context(),username);err!=nil{
		c.Logger().Error("Error updating status of hasread column: ",err)
	}

	count,err:=app.models.Update.GetTotalUnread(c.Request().Context(),username);if err!=nil{
		c.Logger().Error("Error getting total updates: ",err)
	}

	errchan:=app.websocket.PushToClients([]byte(strconv.Itoa(count)),
		func(w WSUser) bool {
		return w.Username==username
	})

	app.CheckChannelError(errchan)
	
	return c.JSON(http.StatusOK, updates)
}