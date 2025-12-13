package main

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/Raaffs/FluxMap/internal/env"
	sessionvar "github.com/Raaffs/FluxMap/internal/sessionVar"
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
	e.DELETE("/api/project/:id/user",app.RemoveUser,app.ManagerLevelAccess)

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
	e.GET("/api/project/:id/tasks", app.GetTasks, IsAuthorizedUser,app.HasProjectAccess)
	e.POST("/api/project/:id/task", app.CreateTask, app.ManagerLevelAccess)
	e.GET("/api/project/:id/task/:taskID", app.GetTaskByID, IsAuthorizedUser)
	e.DELETE("/api/project/:id/task/:taskID",app.RemoveTask,app.ManagerLevelAccess)
	e.PUT("/api/project/:id/task/:taskID/manager", app.ManagerRestrictedTask, app.ManagerLevelAccess)

	e.PUT("/api/project/:id/task/:taskID", app.UpdateUserTask, IsAuthorizedUser)

	e.PUT("/api/project/:id/task/:taskID/approve", app.ApproveTask, app.ManagerLevelAccess)
	e.PUT("/api/project/:id/task/:taskID/assign", app.ManagerRestrictedTask, app.ManagerLevelAccess)

	// PERT & CPM routes
	e.GET("/api/project/:id/pert", app.GetPert, IsAuthorizedUser,app.HasProjectAccess)
	e.POST("/api/project/:id/pert", app.CreatePert, IsAuthorizedUser,app.HasProjectAccess)

	e.GET("/api/project/:id/cpm", app.GetCpm, IsAuthorizedUser,app.HasProjectAccess)
	e.POST("/api/project/:id/cpm", app.CreateCpm, IsAuthorizedUser,app.HasProjectAccess,app.ManagerLevelAccess)

	e.GET("/api/project/:id/update",app.GetProjectUpdates,IsAuthorizedUser,app.HasProjectAccess)
	e.GET("/api/updates",app.GetAllUpdates,IsAuthorizedUser)
	e.GET("/api/updates/recent",app.GetRecentUpdates,IsAuthorizedUser)

	e.GET("/api/ws",app.HandlWS,IsAuthorizedUser)

	//graphs 
	e.GET("/api/project/breakdown",func(c echo.Context) error {
sess, err := session.Get(sessionvar.SESSION_NAME, c)
if err != nil {
    return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get session"})
}

// helper to safely get the length of a slice in session
getRoleCount := func(key string) int {
    if val, ok := sess.Values[key].([]int); ok {
        return len(val)
    }
    	return 0
	}
	roles := []string{string(AdminRole), string(ManagerRole), string(UserRoleVal)}
	total := 0
		for _, role := range roles {
		    total += getRoleCount(role)
		}
	return c.JSON(http.StatusOK, map[string]int{"breakdown": total})	
})
	e.GET("/api/graph/tasks/completed",app.GetTaskCompletedGraph,IsAuthorizedUser)
	e.GET("/api/graph/tasks/approved",app.GetTaskApprovedGraph,IsAuthorizedUser)
	e.GET("/api/graph/tasks/breakdown",app.GetTaskBreakdown,IsAuthorizedUser)
	e.GET("/api/graph/tasks/overdue",app.GetOverDueTasks,IsAuthorizedUser)
}

