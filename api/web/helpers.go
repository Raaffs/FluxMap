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
	"sync"
	"time"

	"slices"

	"github.com/Raaffs/FluxMap/api/external"
	"github.com/Raaffs/FluxMap/internal/models"
	sessionvar "github.com/Raaffs/FluxMap/internal/sessionVar"
	validator "github.com/Raaffs/FluxMap/internal/validators"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
)

type UserRole string

const (
	AdminRole   UserRole = "admin"
	ManagerRole UserRole = "manager"
	UserRoleVal UserRole = "user"
)

type ExistenceStatus int

const (
	Exists ExistenceStatus = iota
	NotExists
	ErrorCheckingExistStatus
)

var (
	NOTIFY_MSG = "msg"
 	NOTIFY_TARGET_USERNAME = "targetUsername"
 	NOTIFY_TARGET_TYPE = "targetType"

	oAuthGoogle="https://www.googleapis.com/oauth2/v2/userinfo"
)

type ProjectResult struct {
	AdminProjects    []*models.Project
	ManagerProjects  []*models.Project
	AssignedProjects []*models.Project
	Err              error
}

var (
	ErrInvalidJson    = errors.New("Invalid JSON")
	ErrFetchingResult = errors.New("Error getting analytics")
)

func FormatDate(t time.Time) string {
	return t.Format("dd-mm-yyyy")
}

func AppendToSessionArray(session *sessions.Session, key string, value int) error {
	// Check if value exists and is of correct type
	arr, ok := session.Values[key].([]int)
	if !ok {
		// Initialize if not present or wrong type
		arr = []int{}
	}
	session.Values[key] = append(arr, value)
	return nil
}

func (app *Application) CacheUserProjectsToSession(c echo.Context) (map[string]any, error) {
	projectschan := make(chan ProjectResult)

	sess, err := session.Get(sessionvar.SESSION_NAME, c)
	if err != nil {
		return nil, err
	}

	username, ok := sess.Values[sessionvar.USERNAME].(string)
	if !ok {
		return nil, c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	go app.FetchProjects(ctx, username, projectschan)
	projects := <-projectschan

	if projects.Err != nil {
		return nil, projects.Err
	}

	projectRoleMap := map[string][]*models.Project{
		string(AdminRole):   projects.AdminProjects,
		string(ManagerRole): projects.ManagerProjects,
		string(UserRoleVal): projects.AssignedProjects,
	}

	for role, projectList := range projectRoleMap {
		for _, project := range projectList {
			err := AppendToSessionArray(sess, role, project.ProjectID)
			if err != nil {
				log.Printf("error appending to session for role %s: %v", role, err)
			}
		}
	}

	roleToProjectMap := make(map[string]any)
	roleToProjectMap[string(AdminRole)] = sess.Values[string(AdminRole)]
	roleToProjectMap[string(ManagerRole)] = sess.Values[string(ManagerRole)]
	roleToProjectMap[string(UserRoleVal)] = sess.Values[string(UserRoleVal)]

	return roleToProjectMap, err
}

func (app *Application) EnsureExists(c echo.Context, checkFunc func(context.Context) (bool, error)) ExistenceStatus {
	ctx := c.Request().Context()
	exist, err := checkFunc(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exist = false
		} else {
			c.Logger().Errorf("Failed checking existence: %v", err)
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			return ErrorCheckingExistStatus
		}
	}
	if !exist {
		return NotExists
	}
	return Exists
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (app *Application) getUserRole(ctx context.Context, username, resourceID string) (UserRole, error) {
	isAdmin, err := app.models.Users.IsAdmin(ctx, username, resourceID)
	if err != nil {
		return "", err
	}
	if isAdmin {
		return AdminRole, nil
	}

	isManager, err := app.models.Users.IsManager(ctx, username, resourceID)
	if err != nil {
		return "", err
	}

	if isManager {
		return ManagerRole, nil
	}
	return UserRoleVal, nil
}

func (app *Application) GetTaskBasedOnAccess(managerFunc echo.HandlerFunc, userFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get(sessionvar.USERNAME).(string)
		// Get the user's role
		projectID := c.Param("id")
		if _, err := strconv.Atoi(projectID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid projectID"})
		}
		role, err := app.getUserRole(c.Request().Context(), username, projectID)
		if err != nil {
			c.Logger().Error("error getting user role: ", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		// Execute based on role
		if role == AdminRole || role == ManagerRole {
			return managerFunc(c)
		}
		return userFunc(c)
	}
}

func SetSession(c echo.Context) (*sessions.Session,error) {
	sess, err := session.Get(sessionvar.SESSION_NAME, c);if err!=nil{
		return nil,err
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	return sess, nil
}

func SaveSession(c echo.Context, sess *sessions.Session, values map[string]any) error {
	for key, value := range values {
		sess.Values[key] = value
	}
	return sess.Save(c.Request(), c.Response())
}

func GetUsernameFromSession(c echo.Context) (string, error) {
	sess, err := session.Get(sessionvar.SESSION_NAME, c)
	if err != nil {
		return "", err
	}
	username, ok := sess.Values[sessionvar.USERNAME].(string)
	if !ok {
		return "", errors.New("username not found in session")
	}
	return username, nil
}

func ValidateManagerUpdate(t models.Task) *validator.Validator {
	v := validator.New()

	v.Check(
		t.AssignedUsername.Valid && t.AssignedUsername.String != "",
		"user",
		"no valid user",
	)
	v.Check(
		t.Approved.Valid,
		"approved",
		"invalid approved status",
	)

	v.Check(
		validator.MinNameLength(t.TaskName),
		validator.ErrNameTooShort.Key,
		validator.ErrNameTooShort.Message,
	)

	v.Check(
		t.TaskDescription.Valid && validator.MinDescriptionLength(t.TaskDescription.String),
		validator.ErrDescriptionTooShort.Key,
		validator.ErrDescriptionTooShort.Message,
	)

	v.Check(
		// TaskStartDate is optional, so it can be null. If a date is provided, it must follow the correct "yyyy-mm-dd" format.
		!t.TaskStartDate.Valid || validator.IsValidDate(t.TaskStartDate.Time.Format("2006-01-02")),
		validator.ErrInvalidDate.Key,
		validator.ErrInvalidDate.Message,
	)

	v.Check(
		!t.TaskDueDate.Valid || validator.IsValidDate(t.TaskDueDate.Time.Format("2006-01-02")),
		validator.ErrInvalidDate.Key,
		validator.ErrInvalidDate.Message,
	)

	v.Check(
		t.TaskStatus.Valid && (t.TaskStatus.String == "pending" || t.TaskStatus.String == "completed"),
		"status",
		"invalid status 2",
	)
	return v
}

func ValidateTask(t models.Task) *validator.Validator {
	v := validator.New()
	v.Check(
		validator.MinNameLength(t.TaskName),
		validator.ErrNameTooShort.Key,
		validator.ErrNameTooShort.Message,
	)

	v.Check(
		validator.MinNameLength(t.TaskDescription.String) && t.TaskDescription.Valid,
		validator.ErrDescriptionTooShort.Key,
		validator.ErrDescriptionTooShort.Message,
	)

	v.Check(
		t.TaskStatus.String == "completed" || t.TaskStatus.String == "pending" && t.TaskStatus.Valid,
		"status",
		"invalid status",
	)

	return v
}

func (app *Application) CheckInvitationStatus(c echo.Context, t models.Task, projectID int, username string) bool {
	//gonna use this hack for now, will think of a better way later
	isAdmin, err := app.models.Users.IsAdmin(c.Request().Context(), username, c.Param("id"))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			c.Logger().Error("error getting access level: ", err)
			c.JSON(http.StatusNotFound, map[string]string{"error": "internal server error"})
			return false
		}
		isAdmin = false
	}
	//check if a	dmin is assigning task to themselves
	if isAdmin && t.AssignedUsername.String == username {
		return true
	}

	switch status := app.EnsureExists(c, func(ctx context.Context) (bool, error) {
		return app.models.Invitation.AcceptedByUser(ctx, projectID, t.AssignedUsername.String)
	}); status {
	case ErrorCheckingExistStatus:
		// error response already sent by EnsureExists, just exit
		return false
	case NotExists:
		c.Logger().Error("user has not accepted invitation or does not exist in this project")
		c.JSON(http.StatusNotFound, map[string]string{"error": "User isn't part of the project or hasn't accepted the invitation yet"})
		return false
	case Exists:
		return true
	}
	return false
}

func (app *Application) GenerateUpdateMessage(ctx context.Context, updatedTask models.Task) (string, error) {
	// Get the old task from DB
	oldTask, err := app.models.Task.GetByID(ctx, updatedTask.TaskID)
	if err != nil {
		return "", err
	}

	var messages []string

	if updatedTask.TaskName != oldTask.TaskName {
		messages = append(messages, fmt.Sprintf("Task name changed from \"%s\" to \"%s\".", oldTask.TaskName, updatedTask.TaskName))
	}

	if updatedTask.TaskDescription != oldTask.TaskDescription {
		messages = append(messages, fmt.Sprintf("Task description for task %s was updated.", updatedTask.TaskName))
	}

	if updatedTask.AssignedUsername.Valid && updatedTask.AssignedUsername.String != oldTask.AssignedUsername.String {
		messages = append(messages, fmt.Sprintf("Assigned user for task %s changed from %s to %s.", updatedTask.TaskName, oldTask.AssignedUsername.String, updatedTask.AssignedUsername.String))
	}

	if updatedTask.TaskDueDate.Valid && !updatedTask.TaskDueDate.Time.Equal(oldTask.TaskDueDate.Time) {
		messages = append(messages, fmt.Sprintf("Due date for task %s changed from %s to %s.",
			updatedTask.TaskName,
			oldTask.TaskDueDate.Time.Format("2006-01-02"),
			updatedTask.TaskDueDate.Time.Format("2006-01-02")))
	}

	if updatedTask.Approved.Valid && updatedTask.Approved.Bool != oldTask.Approved.Bool {
		status := map[bool]string{true: "approved", false: "disapproved"}
		messages = append(messages, fmt.Sprintf("Task %s was %s.", updatedTask.TaskName, status[updatedTask.Approved.Bool]))
	}

	if len(messages) == 0 {
		return "No changes detected.", nil
	}
	return strings.Join(messages, "\n"), nil
}

func (app *Application) FetchProjects(ctx context.Context, username string, resultChan chan<- ProjectResult) {
	var wg sync.WaitGroup
	wg.Add(3)

	adminChan := make(chan []*models.Project, 1)
	managerChan := make(chan []*models.Project, 1)
	assignedChan := make(chan []*models.Project, 1)
	errorChan := make(chan error, 3) // buffered so goroutines don’t block
	done := make(chan struct{})
	go func() {
		defer wg.Done()
		projects, err := app.models.Projects.RetrieveAdminProjects(ctx, username)
		if err != nil {
			errorChan <- err
			return
		}
		adminChan <- projects
	}()

	go func() {
		defer wg.Done()
		projects, err := app.models.Projects.RetrieveManagerProjects(ctx, username)
		if err != nil {
			errorChan <- err
			return
		}
		managerChan <- projects
	}()

	go func() {
		defer wg.Done()
		projects, err := app.models.Projects.RetrieveAssginedProjects(ctx, username)
		if err != nil {
			errorChan <- err
			return
		}
		assignedChan <- projects
	}()

	// Wait for all goroutines
	go func() {
		wg.Wait()
		close(done)
	}()

	var res ProjectResult

	for {
		select {
		case admin, ok := <-adminChan:
			if ok {
				res.AdminProjects = admin
			}
		case manager, ok := <-managerChan:
			if ok {
				res.ManagerProjects = manager
			}
		case assigned, ok := <-assignedChan:
			if ok {
				res.AssignedProjects = assigned
			}
		case err, ok := <-errorChan:
			if ok && err != nil {
				res.Err = err
				resultChan <- res
				return
			}
		case <-done:
			resultChan <- res
			return
		case <-ctx.Done():
			res.Err = ctx.Err()
			resultChan <- res
			return
		}
	}
}

// todo: test this function and replace it with the original fetchprojects function
func (app *Application) FetchProjects2(ctx context.Context, username string, resultChan chan<- ProjectResult) {
	var res ProjectResult
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		projects, err := app.models.Projects.RetrieveAdminProjects(ctx, username)
		if err == nil {
			res.AdminProjects = projects
		}
		return err
	})

	g.Go(func() error {
		projects, err := app.models.Projects.RetrieveManagerProjects(ctx, username)
		if err == nil {
			res.ManagerProjects = projects
		}
		return err
	})

	g.Go(func() error {
		projects, err := app.models.Projects.RetrieveAssginedProjects(ctx, username)
		if err == nil {
			res.AssignedProjects = projects
		}
		return err
	})

	err := g.Wait()
	res.Err = err
	resultChan <- res
}

func (app *Application) UserTasks(c echo.Context) error {
	sess, err := session.Get(sessionvar.SESSION_NAME, c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	username, ok := sess.Values[sessionvar.USERNAME].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "you're not authorized"})
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("not found tasks", id)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "invalid project id"})
	}
	tasks, err := app.models.Task.GetAssignedToUser(c.Request().Context(), id, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, tasks)
}

func (app *Application) ManagerTasks(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("not found tasks manager", id)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "invalid project id"})
	}
	tasks, err := app.models.Task.Get(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, tasks)
}

func detectCycle(graph map[int][]int) error {
	visited := make(map[int]bool)
	onStack := make(map[int]bool)

	var dfs func(int) bool
	dfs = func(node int) bool {
		if onStack[node] {
			return true // Cycle detected
		}
		if visited[node] {
			return false
		}
		visited[node] = true
		onStack[node] = true
		if slices.ContainsFunc(graph[node], dfs) {
			return true
		}
		onStack[node] = false
		return false
	}

	for node := range graph {
		if !visited[node] {
			if dfs(node) {
				return errors.New("dependency cycle detected")
			}
		}
	}
	return nil
}

// Detect cycle for CPM
func DetectCycleCpm(tasks []models.Cpm) error {
	graph := make(map[int][]int)
	for _, task := range tasks {
		for _, dep := range task.Dependencies {
			graph[dep] = append(graph[dep], task.TaskID)
		}
	}
	return detectCycle(graph)
}

func DetectCyclePert(tasks []models.Pert) error {
	graph := make(map[int][]int)
	for _, task := range tasks {
		if task.PredecessorTaskID.Valid {
			graph[int(task.PredecessorTaskID.Int64)] = append(graph[int(task.PredecessorTaskID.Int64)], task.ParentTaskID)
		}
	}
	return detectCycle(graph)
}

func StoreResult[U models.Analytic, T models.ReadDatabase[U]](t T, ctx context.Context, id int, result models.Result) error {
	if err := t.InsertResult(ctx, id, result); err != nil {
		return err
	}
	return nil
}

func Calculate[U models.Analytic, T models.ReadDatabase[U]](v T, ctx context.Context, projectID int) error {
	data, err := v.GetData(ctx, projectID)
	if err != nil {
		return err
	}
	if data == nil {
		return models.ErrRecordNotFound
	}

	result, err := external.RequestAndCalculatePERTCPM(data)
	if err != nil {
		log.Println("Error fetching result: ", err)
		return ErrFetchingResult
	}

	if err := StoreResult(v, ctx, projectID, result); err != nil {
		return err
	}
	return nil
}

func GetAnalytics[U models.Analytic, T models.ReadDatabase[U]](v T, ctx context.Context, id int) ([]*U, models.Result, error) {
	data, err := v.GetData(ctx, id)
	if err != nil {
		return nil, models.Result{}, err
	}

	if data == nil {
		return nil, models.Result{}, models.ErrRecordNotFound
	}

	result, err := v.GetResult(ctx, id)
	if err != nil {
		if !errors.Is(err, models.ErrRecordNotFound) {
			return data, models.Result{}, ErrFetchingResult
		}
	}

	return data, result, nil
}
