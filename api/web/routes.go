package main

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/Raaffs/FluxMap/internal/env"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)
func (app *Application) LoadMiddleware(e *echo.Echo){

	// Middleware setup
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(app.env[env.SESSION_SECRET]))))
	// Registering types for session storage
	gob.Register(map[string][]int{})
	gob.Register(map[string]string{})

	// Rate limiter configuration
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10), Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}
	e.Use(middleware.RateLimiterWithConfig(config))
	// Initialize routes separately
}

func (app *Application)RegisterRoutes(e *echo.Echo){
	// Auth routes
	e.POST("/api/login", app.Login)
	e.POST("/api/register", app.Register)
	e.GET("/api/auth/session",app.SessionCheck)
	e.POST("/api/logout", app.Logout)

	// Project routes
	e.GET("/api/projects", app.GetProjects, IsAuthorizedUser)
	e.POST("/api/projects", app.CreateProject, IsAuthorizedUser)
	e.GET("/api/project/:id", app.GetProjectByID, IsAuthorizedUser)
	e.PUT("/api/project/:id", app.UpdateProject, app.ManagerLevelAccess)

	// Admin & Manager routes
	e.PUT("/api/project/admin/:id", app.UpdateProject, app.AdminLevelAccess)

	e.GET("/api/projects/admin", app.GetAdminProjects, IsAuthorizedUser)
	e.GET("/api/projects/manager", app.GetManagerProjects, IsAuthorizedUser)
	e.GET("/api/projects/assigned", app.GetAssignedProjects, IsAuthorizedUser)

	// Invitation routes
	e.POST("/api/project/:id/invite", app.Invite,IsAuthorizedUser)
	e.GET("/api/invitation", app.GetInvitations,IsAuthorizedUser)
	e.PUT("/api/invitation/:id", app.ConfirmInvitation,IsAuthorizedUser)
	e.GET("/api/invitation/:id/confirmed",app.GetConfirmedUsers,IsAuthorizedUser)
	
	// Task routes
	e.GET("/api/project/:id/tasks", app.GetTasks, IsAuthorizedUser)
	e.POST("/api/project/:id/task", app.CreateTask, app.ManagerLevelAccess, app.SendNotification)
	e.GET("/api/project/:id/task/:taskID", app.GetTaskByID, IsAuthorizedUser)
	e.PUT("/api/project/:id/task/:taskID/manager", app.ManagerRestrictedTask, app.SendNotification, app.ManagerLevelAccess)

	e.PUT("/api/project/:id/task/:taskID", app.UpdateUserTask, IsAuthorizedUser, app.SendNotification)

	e.PUT("/api/project/:id/task/:taskID/approve", app.ApproveTask, app.ManagerLevelAccess)
	e.PUT("/api/project/:id/task/:taskID/assign", app.ManagerRestrictedTask, app.ManagerLevelAccess)

	// PERT & CPM routes
	e.GET("/api/project/:id/pert", app.GetPert, IsAuthorizedUser)
	e.POST("/api/project/:id/pert", app.CreatePert, IsAuthorizedUser)

	e.GET("/api/project/:id/cpm", app.GetCpm, IsAuthorizedUser)
	e.POST("/api/project/:id/cpm", app.CreateCpm, IsAuthorizedUser)

	e.GET("/api/project/:id/update",app.GetProjectUpdates,IsAuthorizedUser)
	e.GET("/api/updates",app.GetAllUpdates,IsAuthorizedUser)

	e.GET("/api/ws",app.HandlWS,IsAuthorizedUser)
}

