package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)
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
