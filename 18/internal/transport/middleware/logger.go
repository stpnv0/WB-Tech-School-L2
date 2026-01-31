package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggingMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set("request_id", reqID)

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		fields := []interface{}{
			"method", c.Request.Method,
			"url", c.Request.URL.String(),
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"request_id", reqID,
			"remote_addr", c.ClientIP(),
		}

		switch {
		case status >= 500:
			log.Error("HTTP request", fields...)
		case status >= 400:
			log.Warn("HTTP request", fields...)
		default:
			log.Info("HTTP request", fields...)
		}
	}
}
