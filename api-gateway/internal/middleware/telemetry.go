package middleware

import (
	"E-Commerce/api-gateway/internal/observability"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Telemetry() gin.HandlerFunc {
	tracer := otel.Tracer("api-gateway")

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Create a span for this request
		ctx, span := tracer.Start(c.Request.Context(), "http_request",
			trace.WithAttributes(
				attribute.String("http.method", method),
				attribute.String("http.path", path),
			),
		)
		defer span.End()

		// Update context with span
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		// Record request duration
		observability.RequestDuration.WithLabelValues(method, path).Observe(duration)

		// Record request count
		observability.RequestCounter.WithLabelValues(
			method,
			path,
			string(rune(status)),
		).Inc()

		// Add response status to span
		span.SetAttributes(attribute.Int("http.status_code", status))
		if status >= 400 {
			span.SetAttributes(attribute.Bool("error", true))
		}
	}
}
