package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sachinggsingh/PDTS/auth-service/internal/utils"
)

// CustomLogger returns a gin middleware that uses the custom logger
func CustomLogger(logger *utils.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get method
		method := c.Request.Method

		// Build query string
		if raw != "" {
			path = path + "?" + raw
		}

		// Log based on status code
		msg := fmt.Sprintf("%s | %13v | %15s | %-7s %s",
			statusColor(statusCode),
			latency,
			clientIP,
			method,
			path,
		)

		if statusCode >= 500 {
			logger.Error(msg)
		} else if statusCode >= 400 {
			logger.Warn(msg)
		} else {
			logger.Info(msg)
		}
	}
}

// statusColor returns a color code based on HTTP status code
func statusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return fmt.Sprintf("\033[32m%3d\033[0m", code) // Green
	case code >= 300 && code < 400:
		return fmt.Sprintf("\033[34m%3d\033[0m", code) // Blue
	case code >= 400 && code < 500:
		return fmt.Sprintf("\033[33m%3d\033[0m", code) // Yellow
	default:
		return fmt.Sprintf("\033[31m%3d\033[0m", code) // Red
	}
}
