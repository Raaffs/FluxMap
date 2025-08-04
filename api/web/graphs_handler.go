package main

import (
	"database/sql"
	"errors"
	"net/http"

	sessionvar "github.com/Raaffs/FluxMap/internal/sessionVar"
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