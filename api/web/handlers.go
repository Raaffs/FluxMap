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
	defer app.CacheUserProjectsToSession(c)
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
	log.Println("sess vals login", sess.Values)
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		c.Logger().Error("error saving session: ", err)
		return c.JSON(http.StatusInternalServerError, "error saving session")
	}

	return c.JSON(http.StatusOK, "")
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
			return c.JSON(http.StatusConflict, MapMessage("error", "User already exist"))
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
		log.Println("Error json: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON payload222",
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
		c.JSON(http.StatusInternalServerError, MapMessage("Project", "An error occurred while retrieving project"))
	}
	c.JSON(http.StatusOK, adminProjects)
	return nil
}

func (app *Application) GetManagerProjects(c echo.Context) error {
	username := c.Get(sessionvar.USERNAME).(string)

	managerProjects, err := app.models.Projects.RetrieveManagerProjects(c.Request().Context(), username)
	if err != nil {
		c.Logger().Error("Error retrieving manager projects: ", err)
		return c.JSON(http.StatusInternalServerError, MapMessage("Project", "An error occurred while retrieving manager projects"))
	}

	return c.JSON(http.StatusOK, managerProjects)
}

func (app *Application) GetAssignedProjects(c echo.Context) error {
	username := c.Get(sessionvar.USERNAME).(string)

	assignedProjects, err := app.models.Projects.RetrieveAssginedProjects(c.Request().Context(), username)
	if err != nil {
		// Log and handle errors
		c.Logger().Error("Error retrieving assigned projects: ", err)
		return c.JSON(http.StatusInternalServerError, MapMessage("Project", "An error occurred while retrieving assigned projects"))
	}
	return c.JSON(http.StatusOK, assignedProjects)
}

func (app *Application) GetProjectByID(c echo.Context) error {
	id := c.Param("id")
	projID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, MapMessage("message", "Invalid project id"))
	}


	projects, err := app.models.Projects.RetrieveProjectByID(c.Request().Context(), projID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, MapMessage("message", "Project not found"))
		}
		c.Logger().Error("Error retrieving project by id: ", err)
		return c.JSON(http.StatusInternalServerError, MapMessage("message", "An error occurred while retrieving project by id"))
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

	exist, err := app.models.Users.Exist(c.Request().Context(), invitation.Username)
	if !exist {
		c.Logger().Error("user doesn't exists")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user doesn't exist"})
	}

	if err != nil {
		c.Logger().Error("Error inviting user:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	alreadyInvited, err := app.models.Invitation.Exist(c.Request().Context(), invitation.Username, id)
	if err != nil {
		c.Logger().Error("Error checking invitation:", err)
		return c.JSON(http.StatusConflict, map[string]string{"error": "internal server error"})
	}
	if alreadyInvited {
		return c.JSON(http.StatusConflict, map[string]string{"error": "user is already invited"})
	}

	if err := app.models.Invitation.Invite(c.Request().Context(), invitation.Username, id, invitation.Role); err != nil {
		c.Logger().Error("Error inviting user: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return nil
}

func (app *Application) GetInvitations(c echo.Context) error {
	username := c.Get(sessionvar.USERNAME).(string)

	invitations, err := app.models.Invitation.GetPendingInvitations(c.Request().Context(), username)
	if err != nil {

		c.Logger().Error("Error getting invitations: ", err)
		return c.JSON(http.StatusInternalServerError, MapMessage("Invitations", "Failed to get invitations"))
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

func (app *Application) CreateTask(c echo.Context) error {
	var t models.Task
	username := c.Get(sessionvar.USERNAME).(string)
	v := validator.New()
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Json payload"})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Warn(MapMessage("error converting to string", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}
	//gonna use this hack for now, will think of a better way later
	isAdmin, ok := c.Get("isAdmin").(bool)
	if !ok {
		isAdmin = false
	}
	//check if they're part of project
	accepted, err := app.models.Invitation.HasAcceptedInvitation(c.Request().Context(), id, t.AssignedUsername.String)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Logger().Warn("user is has not accepted invitation")

			return c.JSON(http.StatusForbidden, map[string]string{"error": "User isn't part of the project or hasn't accepted the invitation yet"})
		}
		c.Logger().Error("error getting invitation status: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	if !accepted && !isAdmin {
		c.Logger().Warn("user is has not accepted invitation")
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User isn't part of the project or hasn't accepted the invitation yet"})
	}

	t.ParentProjectID = id
	v.Check(
		validator.MinNameLength(t.TaskName),
		validator.ErrNameTooShort.Key,
		validator.ErrNameTooShort.Message,
	)
	if !v.Valid() {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": v.Errors})
	}

	t.Createdby = username

	if err := app.models.Task.Create(c.Request().Context(), t); err != nil {
		c.Logger().Error(MapMessage("Error creating task", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "An error occured"})
	}
	u := models.Update{
		Msg:            fmt.Sprintf("Task %s created by %s assigned to %s", t.TaskName, username, t.AssignedUsername.String),
		TargetType:     "user",
		TargetUsername: t.AssignedUsername,
	}
	SetNotifyContext(c, u.Msg, u.TargetType, u.TargetUsername.String)
	return c.JSON(http.StatusOK, "task created")
}

func (app *Application) GetTasks(c echo.Context) error {
	taskfunc := app.GetTaskBasedOnAccess(app.ManagerTasks, app.UserTasks)
	return taskfunc(c)
}

func (app *Application) GetTaskByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.Logger().Error(MapMessage("error converting to string", err.Error()))
		return c.JSON(http.StatusNotFound, MapMessage("error", "task not found"))
	}
	task, err := app.models.Task.GetTaskByID(c.Request().Context(), id)
	if err != nil {
		c.Logger().Error(MapMessage("Error retrieving task", err.Error()))
		return c.JSON(http.StatusInternalServerError, MapMessage("error", "failed to retrieve task"))
	}

	return c.JSON(http.StatusOK, task)
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
	username := c.Get(sessionvar.USERNAME).(string)
	v := validator.New()
	id, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		return c.JSON(http.StatusNotFound, MapMessage("error", "task not found"))
	}
	if err := c.Bind(&status); err != nil {
		return c.JSON(http.StatusBadRequest, MapMessage("error", "invalid request"))
	}
	v.Check(
		status.TaskStatus == "completed" || status.TaskStatus == "pending",
		"status",
		"invalid status",
	)
	taskName, err := app.models.Task.UpdateTask(c.Request().Context(), id, status.TaskStatus)
	if err != nil {
		c.Logger().Error("error updating task: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	u := models.Update{
		Msg: fmt.Sprintf(
			"Task %s has been set to %s by %s",
			taskName,
			status.TaskStatus,
			username,
		),
		TargetType: "all",
	}
	SetNotifyContext(c, u.Msg, u.TargetType, u.TargetUsername.String)
	return nil
}

func (app *Application) UpdateManagerTask(c echo.Context) error {
	var t models.Task
	v := validator.New()

	if err := c.Bind(&t); err != nil {
		c.Logger().Warn("Failed to bind request body: ", err)
		return c.JSON(http.StatusBadRequest, MapMessage("error", "Invalid request"))
	}

	v.Check(
		validator.MinNameLength(t.TaskName),
		validator.ErrNameTooShort.Key,
		validator.ErrNameTooShort.Message,
	)

	if !v.Valid() {
		c.Logger().Warnf("Validation failed for task: %v", v.Errors)
		return c.JSON(http.StatusBadRequest, v)
	}

	if !t.AssignedUsername.Valid {
		c.Logger().Warn("Task submission missing assigned user")
		return c.JSON(http.StatusBadRequest, MapMessage("error", "Task must be assigned to a user"))
	}

	if t.Approved.Bool {
		t.TaskApprovedDate.Time = time.Now()
	}

	if err := app.models.Task.UpdateManagerTask(c.Request().Context(), t); err != nil {
		c.Logger().Error("Error updating task in DB: ", err)
		return c.JSON(http.StatusInternalServerError, MapMessage("error", "Internal server error"))
	}

	return c.JSON(http.StatusOK, "Updated successfully")
}

func (app *Application) ApproveTask(c echo.Context) error {
	var t models.Task
	username := c.Get(sessionvar.USERNAME).(string)
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Warn("Error converting project id from string to int : ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}
	id := c.Param("taskID")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Invalid task ID")
	}

	if err := c.Bind(&t); err != nil {
		c.Logger().Warn("Error reading json: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	t.TaskID = taskID

	accepted, err := app.models.Invitation.HasAcceptedInvitation(c.Request().Context(), projectID, t.AssignedUsername.String)
	if err != nil {
		c.Logger().Error("Error getting invitation status: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	if !accepted {
		c.Logger().Warn("User %s has not accepted invitation to project id %d", t.AssignedUsername.String, projectID)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User has not accepted invitation to project"})
	}

	v := ValidateManagerUpdate(t)

	if !v.Valid() {
		c.Logger().Warnf("Validation failed for task: %v", v.Errors)
		return c.JSON(http.StatusBadRequest, v)
	}
	if t.Approved.ValueOrZero() {
		t.TaskApprovedDate = null.NewTime(time.Now(), true)
	} else {
		t.TaskApprovedDate = null.NewTime(time.Time{}, false)
	}

	u := models.Update{
		Msg: fmt.Sprintf(
			"Task %s assigned to %s %s by %s",
			t.TaskName,
			t.AssignedUsername.String,
			map[bool]string{true: "approved", false: "disapproved"}[t.Approved.ValueOrZero()],
			username,
		),
		TargetType:     "user",
		TargetUsername: t.AssignedUsername,
	}
	
	if err := app.models.Task.UpdateManagerTask(c.Request().Context(), t); err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			c.Logger().Warn("Task not found :", err)
			return c.JSON(http.StatusNotFound, MapMessage("message", models.ErrRecordNotFound.Error()))
		}
		c.Logger().Error("error updating manager task : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update task"})
	}

	SetNotifyContext(c, u.Msg, u.TargetType, u.TargetUsername.String)
	return c.JSON(http.StatusOK, map[string]string{"message": "task approved successfully"})
}

func (app *Application) ManagerRestrictedTask(c echo.Context) error {
	var t models.Task
	username := c.Get(sessionvar.USERNAME).(string);
	if err := c.Bind(&t); err != nil {
		c.Logger().Error("error reading json,", err)
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}
	
	log.Println(username)
	id := c.Param("taskID");taskID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Invalid task ID")
	}
	t.TaskID = taskID
	v:=ValidateManagerUpdate(t)
	if !v.Valid() {
		c.Logger().Error(v)
		return c.JSON(http.StatusBadRequest, v)
	}
	if err := app.models.Task.UpdateManagerTask(c.Request().Context(), t); err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			c.Logger().Warn("Task not found :", err)
			return c.JSON(http.StatusNotFound, MapMessage("message", models.ErrRecordNotFound.Error()))
		}
		c.Logger().Error("error updating manager task : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update task"})
	}
	msg,err:=app.GenerateUpdateMessage(c.Request().Context(),t);if err!=nil{
		c.Logger().Error("error generating update message : ", err)
	}

	u := models.Update{
		Msg: msg,
		TargetType:     "user",
		TargetUsername: t.AssignedUsername,
	}

	SetNotifyContext(c, u.Msg, u.TargetType, u.TargetUsername.String)

	return c.JSON(http.StatusOK, "task approved")
}

func (app *Application) GetPert(c echo.Context) error {
	r := struct {
		Data   []*models.Pert `json:"data"`
		Result map[string]any `json:"result"`
	}{}

	id := c.Param("id")
	projectID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Invalid project ID")
	}

	data, result, err := GetAnalytics(&app.models.Pert, c.Request().Context(), projectID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, MapMessage("message", "No data CPM related data found"))
		}
		if errors.Is(err, ErrFetchingResult) {
			r.Data = data
			return c.JSON(http.StatusPartialContent, r.Data)
		}
		c.Logger().Error("error getting cpm values : ", err)
		return c.JSON(http.StatusInternalServerError, MapMessage("error", "Error getting CPM data"))
	}
	r.Data = data
	r.Result = result.Result

	return c.JSON(http.StatusOK, r)
}

func (app *Application) CreatePert(c echo.Context) error {
	var pert []models.Pert
	if err := c.Bind(&pert); err != nil {
		c.Logger().Warn("error binding pert : ", err)
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
		return c.JSON(http.StatusInternalServerError, MapMessage("Error", "Failed to calculate values"))
	}
	return c.JSON(http.StatusOK, MapMessage("PERT", "data and result inserted successfully"))
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
			return c.JSON(http.StatusNotFound, MapMessage("message", "No data CPM related data found"))
		}
		if errors.Is(err, ErrFetchingResult) {
			r.Data = data
			return c.JSON(http.StatusPartialContent, r.Data)
		}
		c.Logger().Error("error getting cpm values : ", err)
		return c.JSON(http.StatusInternalServerError, MapMessage("error", "Error getting CPM data"))
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
	updates, err := app.models.Update.GetAllUpdates(c.Request().Context(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Logger().Warn("not updates : ", err)
			return c.JSON(http.StatusNotFound, MapMessage("error", "No updates found"))
		}
		c.Logger().Error("error getting updates : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, updates)
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
			return c.JSON(http.StatusNotFound, MapMessage("error", "No updates found"))
		}
		c.Logger().Error("error getting updates : ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, updates)
}
