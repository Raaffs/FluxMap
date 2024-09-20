package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)
func (app *Application)InitRoutes()*echo.Echo{
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	
	e.GET("/",app.Home )
	e.POST("/api/login",app.Login)
	e.POST("/api/register",app.Register)
	e.POST("/api/logout",app.Logout)

	e.GET("/api/projects",app.GetProjects,IsAuthorizedUser)
	e.POST("/api/project",app.CreateProject,IsManager)
	e.GET("/api/project/:id",app.GetProjectByID,IsAuthorizedUser)
	e.PUT("/api/project/:id",app.UpdateProject,IsManager)
	
	e.GET("/api/project/:id/tasks",app.GetTasks,IsAuthorizedUser)
	e.POST("/api/project/:id/task",app.CreateProject,IsManager)
	e.GET("/api/project/:id/task/:id",app.GetTaskByID,IsAuthorizedUser)
	e.PUT("/api/project/:id/task/:id",app.UpdateTask,IsAuthorizedUser)

	e.PUT("/api/project/:id/task/:id/approve",app.ManagerRestrictedTaskUpdate,IsManager)
	e.PUT("/api/project/:id/task/:id/assign",app.ManagerRestrictedTaskUpdate,IsManager)
	
	e.GET("/api/project/:id/pert",app.GetPert,IsAuthorizedUser)
	e.POST("/api/project/:id/pert",app.CreateCpm,IsAuthorizedUser)
	e.PUT("/api/project/:id/pert",app.UpdatePert,IsAuthorizedUser)

	e.GET("/api/project/:id/cpm",app.GetCpm,IsAuthorizedUser)
	e.POST("/api/project/:id/cpm",app.CreateCpm,IsAuthorizedUser)
	e.PUT("/api/project/:id/cpm",app.UpdateCpm,IsAuthorizedUser)
	return e
}