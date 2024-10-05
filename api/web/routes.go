package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)
func (app *Application)InitRoutes()*echo.Echo{
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	// e.Use(middleware.CSRF())
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10), Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string,err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}
	e.Use(middleware.RateLimiterWithConfig(config))
	e.POST("/api/login",app.Login)
	e.POST("/api/register",app.Register)
	e.POST("/api/logout",app.Logout)

	e.GET("/api/project",app.GetProjects,IsAuthorizedUser)
	e.POST("/api/project",app.CreateProject,IsAuthorizedUser)
	e.PUT("/api/project/:id",app.UpdateProject,app.ManagerLevelAccess)
	e.POST("/api/project/:id/invite",app.Invite)
	e.PUT("/api/project/:id/invite",app.Invite)
	e.PUT("/api/project/admin/:id",app.UpdateProject,app.AdminLevelAccess)
	e.POST("/api/project/:id/manager", app.AddManager,app.AdminLevelAccess)
	
	e.GET("/api/project/:id/tasks",app.GetTasks,IsAuthorizedUser)
	e.POST("/api/project/:id/task",app.CreateTask,app.ManagerLevelAccess)
	e.GET("/api/project/:id/task/:taskID",app.GetTaskByID,IsAuthorizedUser)
	e.PUT("/api/project/:id/task/:taskID/manager",app.ManagerRestrictedTask,app.ManagerLevelAccess)
	e.PUT("/api/project/:id/task/:taskID",app.UpdateTask,IsAuthorizedUser)


	e.PUT("/api/project/:id/task/:taskID/approve",app.ManagerRestrictedTask,app.ManagerLevelAccess)
	e.PUT("/api/project/:id/task/:taskID/assign",app.ManagerRestrictedTask,app.ManagerLevelAccess)
	
	e.GET("/api/project/:id/pert",app.GetPert,IsAuthorizedUser)
	e.POST("/api/project/:id/pert",app.CreatePert,IsAuthorizedUser)

	e.GET("/api/project/:id/cpm",app.GetCpm,IsAuthorizedUser)
	e.POST("/api/project/:id/cpm",app.CreateCpm,IsAuthorizedUser)
	return e
}