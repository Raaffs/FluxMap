package main

import (
	"net/http"
	"time"
    "encoding/gob"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
    "github.com/gorilla/sessions"
    "github.com/labstack/echo-contrib/session"
)
func (app *Application)InitRoutes()*echo.Echo{
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization },
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
    gob.Register(map[string][]int{})
	gob.Register(map[string]string{})
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

	e.GET("/api/projects",app.GetProjects,IsAuthorizedUser)
	e.POST("/api/projects",app.CreateProject,IsAuthorizedUser)
	e.GET("/api/project/:id",app.GetProjectByID,IsAuthorizedUser)
	e.PUT("/api/project/:id",app.UpdateProject,app.ManagerLevelAccess)
	e.PUT("/api/project/admin/:id",app.UpdateProject,app.AdminLevelAccess)
	e.POST("/api/project/:id/manager", app.AddManager,app.AdminLevelAccess)

	e.GET("/api/projects/admin",app.GetAdminProjects,IsAuthorizedUser)
	e.GET("/api/projects/manager",app.GetManagerProjects,IsAuthorizedUser)
	e.GET("/api/projects/assigned",app.GetAssignedProjects,IsAuthorizedUser)

	e.POST("/api/project/:id/invite",app.Invite)
	e.GET("/api/invitation",app.GetInvitations)
	e.PUT("/api/invitation",app.ConfirmInvitation)
	
	e.GET("/api/project/:id/tasks",app.GetTasks,IsAuthorizedUser)
	e.POST("/api/project/:id/task",app.CreateTask,app.ManagerLevelAccess)
	e.GET("/api/project/:id/task/:taskID",app.GetTaskByID,IsAuthorizedUser)
	e.PUT("/api/project/:id/task/:taskID/manager",app.ManagerRestrictedTask,app.ManagerLevelAccess)
	e.PUT("/api/project/:id/task/:taskID",app.UpdateUserTask,IsAuthorizedUser)

	e.PUT("/api/project/:id/task/:taskID/approve",app.ManagerRestrictedTask,app.ManagerLevelAccess)
	e.PUT("/api/project/:id/task/:taskID/assign",app.ManagerRestrictedTask,app.ManagerLevelAccess)
	
	e.GET("/api/project/:id/pert",app.GetPert,IsAuthorizedUser)
	e.POST("/api/project/:id/pert",app.CreatePert,IsAuthorizedUser)

	e.GET("/api/project/:id/cpm",app.GetCpm,IsAuthorizedUser)
	e.POST("/api/project/:id/cpm",app.CreateCpm,IsAuthorizedUser)
	return e
}