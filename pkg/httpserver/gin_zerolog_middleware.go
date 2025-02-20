package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

func GinRequestZeroLogger(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract traceparent from context using otel package
		traceID := getTraceIDFromContext(c.Request.Context())

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		stop := time.Since(start)
		status := c.Writer.Status()

		var (
			zevent *zerolog.Event
			msg    string
		)

		switch {
		case status >= http.StatusInternalServerError:
			zevent = logger.
				Error().
				Stack().
				Err(c.Errors.Last())
			msg = "http request failed"
		case status >= http.StatusBadRequest:
			zevent = logger.
				Warn().
				Err(c.Errors.Last())
			msg = "http request failed"
		default:
			zevent = logger.
				Debug()
			msg = "http request received"
		}

		zevent.
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Str("trace_id", traceID).
			Dur("latency_ms", stop).
			Msg(msg)
	}
}

func getTraceIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		traceID := spanCtx.TraceID()
		return traceID.String()
	}
	return ""
}
