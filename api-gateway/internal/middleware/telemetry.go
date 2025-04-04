package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"log"
)

// Telemetry middleware captures request duration and logs it
func Telemetry() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record the start time of the request
		start := time.Now()

		// Process the request
		c.Next()

		// Calculate the duration
		duration := time.Since(start)

		// Log telemetry data (method, path, status, and duration)
		log.Printf(
			"Telemetry: %s %s - Status: %d - Duration: %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}