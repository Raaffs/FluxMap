package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Raaffs/FluxMap/api/external"
	"github.com/Raaffs/FluxMap/internal/models"
	sessionvar "github.com/Raaffs/FluxMap/internal/sessionVar"
	validator "github.com/Raaffs/FluxMap/internal/validators"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	AdminRole   UserRole = "admin"
	ManagerRole UserRole = "manager"
	UserRoleVal UserRole = "user"
)
var NOTIFY_MSG="msg"
var NOTIFY_TARGET_USERNAME="targetUsername"
var NOTIFY_TARGET_TYPE="targetType"

type ProjectResult struct {
	AdminProjects    []*models.Project
	ManagerProjects  []*models.Project
	AssignedProjects []*models.Project
	Err              error
}


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

func SetNotifyContext(c echo.Context, msg,targetType,targetUsername string){
	c.Set(NOTIFY_MSG,msg)
	c.Set(NOTIFY_TARGET_TYPE,targetType)
	c.Set(NOTIFY_TARGET_USERNAME,targetUsername)
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

func (app *Application) CacheUserProjectsToSession(c echo.Context) error {
	projectschan := make(chan ProjectResult)

	sess, err := session.Get(sessionvar.SESSION_NAME, c)
	if err != nil {
		log.Println("sess in store session projects", sess.Values)
		return err
	}

	username, ok := sess.Values[sessionvar.USERNAME].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	ctx,cancel:=context.WithTimeout(c.Request().Context(),10*time.Second)
	defer cancel()

	go app.FetchProjects(ctx,username, projectschan)
	projects := <-projectschan

	if projects.Err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to fetch projects")
	}

	projectRoleMap := map[string][]*models.Project{
		string(AdminRole):   projects.AdminProjects,
		string(ManagerRole): projects.ManagerProjects,
		string(UserRoleVal):    projects.AssignedProjects,
	}

	for role, projectList := range projectRoleMap {
		for _, project := range projectList {
			err := AppendToSessionArray(sess, role, project.ProjectID)
			if err != nil {
				log.Printf("error appending to session for role %s: %v", role, err)
			}
		}
	}
	log.Println(sess.Values)
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
		username:=c.Get(sessionvar.USERNAME).(string);
        // Get the user's role
		id:=c.Param("id")
		if _,err:=strconv.Atoi(id);err!=nil{
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}
        role, err := app.getUserRole(c.Request().Context(), username, id)
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

func ValidateManagerUpdate(t models.Task)(*validator.Validator){
	v:=validator.New()
	
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
		t.TaskStatus.Valid && (t.TaskStatus.String=="pending" || t.TaskStatus.String=="accepted"),	
		"status",
		"invalid status",
	)
	return v
}

func (app *Application) GenerateUpdateMessage(ctx context.Context, updatedTask models.Task) (string, error) {
	// Get the old task from DB
	oldTask, err := app.models.Task.GetTaskByID(ctx, updatedTask.TaskID)
	if err != nil {
		return "", err
	}

	var messages []string

	if updatedTask.TaskName != oldTask.TaskName {
		messages = append(messages, fmt.Sprintf("Task name changed from \"%s\" to \"%s\"", oldTask.TaskName, updatedTask.TaskName))
	}

	if updatedTask.TaskDescription != oldTask.TaskDescription {
		messages = append(messages, "Task description for task %s was updated",updatedTask.TaskName)
	}

	if updatedTask.AssignedUsername.Valid && updatedTask.AssignedUsername.String != oldTask.AssignedUsername.String {
		messages = append(messages, fmt.Sprintf("Assigned user for task %s changed from %s to %s",updatedTask.TaskName, oldTask.AssignedUsername.String, updatedTask.AssignedUsername.String))
	}

	if updatedTask.TaskDueDate.Valid && !updatedTask.TaskDueDate.Time.Equal(oldTask.TaskDueDate.Time) {
		messages = append(messages, fmt.Sprintf("Due date for task %s changed from %s to %s",
			updatedTask.TaskName,
			oldTask.TaskDueDate.Time.Format("2006-01-02"), 
			updatedTask.TaskDueDate.Time.Format("2006-01-02")))
	}

	if updatedTask.Approved.Valid && updatedTask.Approved.Bool != oldTask.Approved.Bool {
		status := map[bool]string{true: "approved", false: "disapproved"}
		messages = append(messages, fmt.Sprintf("Task %s was %s", updatedTask.TaskName,status[updatedTask.Approved.Bool]))
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
	done:=make(chan struct{})
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


func (app *Application)UserTasks(c echo.Context)error{
    sess,err:=session.Get(sessionvar.SESSION_NAME,c);if err!=nil{
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
    }
    log.Println("sess in access task2",sess.Values)
    username,ok:=sess.Values[sessionvar.USERNAME].(string);if !ok{
        return c.JSON(http.StatusUnauthorized,map[string]string{"error":"you're not authorized"})
    }
    id,err:=strconv.Atoi(c.Param("id"));if err!=nil{
        log.Println("not found tasks",id)
        return c.JSON(http.StatusNotFound,map[string]string{"error":"invalid project id"})
    }
    tasks,err:=app.models.Task.GetUserTasks(c.Request().Context(),id,username)
    if err!=nil{
        return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
    }
    return c.JSON(http.StatusOK, tasks)
}

func (app *Application)ManagerTasks(c echo.Context)error{
    id,err:=strconv.Atoi(c.Param("id"));if err!=nil{
        log.Println("not found tasks manager",id)
        return c.JSON(http.StatusNotFound,map[string]string{"error":"invalid project id"})
    }
    tasks,err:=app.models.Task.GetTasks(c.Request().Context(),id)
    if err!=nil{
        return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
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
		for _, neighbor := range graph[node] {
			if dfs(neighbor) {
				return true
			}
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


func StoreResult[U models.Analytic, T models.ReadDatabase[U]](t T,ctx context.Context,id int, result models.Result)error{
	if err:=t.InsertResult(ctx,id,result);err!=nil{
		return err
	}
	return nil
}

func Calculate[U models.Analytic,T models.ReadDatabase[U]](v T,ctx context.Context, id int)(error){
	data,err:=v.GetData(ctx,id);if err!=nil{
		return err 
	}
	if data==nil{
		return models.ErrRecordNotFound
	}

	result,err:=external.RequestAndCalculatePERTCPM(data); if err!=nil{
		log.Println("Error fetching result: ",err)
		return ErrFetchingResult
	}

	if err:=StoreResult(v,ctx,id,result);err!=nil{
		return err
	}
	return nil
}

func GetAnalytics[U models.Analytic,T models.ReadDatabase[U]](v T,ctx context.Context, id int)([]*U,models.Result,error){
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
	return data,result,nil
}
