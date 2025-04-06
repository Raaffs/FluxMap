package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Raaffs/FluxMap/api/external"
	"github.com/Raaffs/FluxMap/internal/models"
	sessionvar "github.com/Raaffs/FluxMap/internal/sessionVar"
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
        sess, err := session.Get(sessionvar.SESSION_NAME, c)
        if err != nil {
            log.Println("sess in access task0", sess.Values)
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing role cookie"})
        }
        log.Println("sess in access task1", sess.Values)

        username, ok := sess.Values[sessionvar.USERNAME].(string)
        if !ok {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "you're not authorized"})
        }

        // Get the user's role
        role, err := app.getUserRole(c.Request().Context(), username, c.Param("id"))
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
