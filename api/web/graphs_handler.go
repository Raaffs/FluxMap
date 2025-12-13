package main

import (
	"database/sql"
	"errors"
	"net/http"

	sessionvar "github.com/Raaffs/FluxMap/internal/sessionVar"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (app *Application)GetTaskCompletedGraph(c echo.Context)error{
	username:=c.Get(sessionvar.USERNAME).(string)
	graph,err:=app.graphs.TaskCompletedByDate(c.Request().Context(),username);if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return c.JSON(http.StatusNotFound,map[string]string{"graph":"no data"})
		}
		app.logger.Error("Error retrieving graph data: ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
	}
	return c.JSON(http.StatusOK,map[string]any{"graph":graph})
}

func (app *Application)GetTaskApprovedGraph(c echo.Context)error{
	username:=c.Get(sessionvar.USERNAME).(string)
	graph,err:=app.graphs.TaskApprovedByDate(c.Request().Context(),username);if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return c.JSON(http.StatusNotFound,map[string]string{"graph":"no data"})
		}
		app.logger.Error("Error retrieving graph data: ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
	}
	return c.JSON(http.StatusOK,map[string]any{"graph":graph})
}

//active pending, overdue, completed
func (app *Application)GetTaskBreakdown(c echo.Context)error{
	username:=c.Get(sessionvar.USERNAME).(string)
	breakdown,err:=app.graphs.TaskStatusBreakdown(c.Request().Context(),username); if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return c.JSON(http.StatusOK,map[string]any{"graph":[3]int{0,0,0}})
		}
		app.logger.Error("Error getting status breakdown graph: ",err)
		return c.JSON(http.StatusInternalServerError,map[string]string{"error":"internal server error"})
	}
	return c.JSON(http.StatusOK,map[string]any{"graph":breakdown})
}

func (app *Application) GetProjectBreakdown(c echo.Context) error {
	sess, err := session.Get(sessionvar.SESSION_NAME, c)
	if err != nil {
		app.logger.Error("Error getting session: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}

	getProjects := func(key string) []int {
		if v, ok := sess.Values[key].([]int); ok {
			return v
		}
		app.logger.Errorf("Session value for %s is missing or invalid", key)
		return []int{}
	}

	totalAdmin := getProjects(sessionvar.ADMIN)
	totalManager := getProjects(sessionvar.MANAGER)
	totalUser := getProjects(sessionvar.USER)

	return c.JSON(http.StatusOK, map[string][]int{
		"projects": {len(totalAdmin), len(totalManager), len(totalUser)},
	})
}


