package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Raaffs/FluxMap/internal/env"
	"github.com/Raaffs/FluxMap/internal/graphs"
	"github.com/Raaffs/FluxMap/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)	

type Application struct {
	env       map[string]string
	models    models.Models
	graphs    graphs.GraphModel
	websocket *ConnectionManager
	logger    echo.Logger
}

type SlogHandler struct {
	slog.Handler
	w io.Writer
}

func (h *SlogHandler) Handle(ctx context.Context, r slog.Record) error {
	fields := make(map[string]any, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	output := map[string]any{
		"time":    r.Time.Format(time.RFC3339Nano),
		"level":   r.Level.String(),
		"message": r.Message,
		"data":    fields,
	}

	// Pretty-print JSON with 2-space indentation
	b, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	// ANSI Color Coding: Red if level is ERROR
	out := string(b)
	if r.Level >= slog.LevelError {
		out = fmt.Sprintf("\033[31m%s\033[0m", out) // Red text
	}

	fmt.Fprintln(h.w, out)
	return nil
}

// StructuredLogger is the Echo Middleware
func StructuredLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process request
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			stop := time.Now()
			req := c.Request()
			res := c.Response()

			// Log Attributes
			attrs := []slog.Attr{
				slog.String("remote_ip", c.RealIP()),
				slog.String("host", req.Host),
				slog.String("method", req.Method),
				slog.String("uri", req.RequestURI),
				slog.String("user_agent", req.UserAgent()),
				slog.Int("status", res.Status),
				slog.Int64("latency", int64(stop.Sub(start))),
				slog.String("latency_human", stop.Sub(start).String()),
				slog.Int64("bytes_in", req.ContentLength),
				slog.Int64("bytes_out", res.Size),
			}

			lvl := slog.LevelInfo
			msg := "request processed"

			if err != nil {
				lvl = slog.LevelError
				msg = err.Error()
				attrs = append(attrs, slog.String("error_detail", fmt.Sprintf("%+v", err)))
			} else if res.Status >= 400 {
				lvl = slog.LevelError
			}
			logger.LogAttrs(context.Background(), lvl, msg, attrs...)

			return nil
		}
	}
}

func connectWithRetry(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	for i := range 10 {
		log.Printf("Attempting DB connection (attempt %d/10)...", i+1)
		
		pool, err = pgxpool.New(ctx, dbURL)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				log.Println("PING SUCCESSFUL")
				return pool, nil
			}
		}

		log.Printf("DB not ready: %v. Retrying in 2s...", err)
		if pool != nil {
			pool.Close()
		}
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to DB: %w", err)
}

func main() {

	envMap := map[string]string{
		env.API_PORT:      	os.Getenv(env.API_PORT),
		env.DB_URL:        	os.Getenv(env.DB_URL),
		env.CLIENT_PORT:   	os.Getenv(env.CLIENT_PORT),
		env.SESSION_SECRET:	os.Getenv(env.SESSION_SECRET),
	}

	// Hard fails for required vars
	if envMap[env.DB_URL] == "" {
		log.Fatal("DB_URL is missing")
	}
	if envMap[env.API_PORT] == "" {
		log.Fatal("API_PORT is missing")
	}

	ctx := context.Background()
	
	conn, err := connectWithRetry(ctx,envMap[env.DB_URL]);if err!=nil{
		log.Fatalf("Could not connect to DB: %v", err)
	}

	app := &Application{
		env:       envMap,
		models:    models.NewModels(conn),
		graphs:    graphs.NewGraphs(conn),
		websocket: NewConnectionManager(),
		logger:    echo.New().Logger,
	}

	router := echo.New()

	app.LoadMiddleware(router)
	app.RegisterRoutes(router)

	port := fmt.Sprintf(":%s", envMap[env.API_PORT])
	log.Println("API starting on", port)

	if err := router.Start(port); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}

