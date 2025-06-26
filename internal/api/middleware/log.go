package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/James-D-Wood/blog-api/internal/constant"
)

// LoggerMiddleware logs each request and adds logger to context
func LoggerMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("receiving request", "method", r.Method, "path", r.URL.Path)
		ctx := context.WithValue(r.Context(), constant.LoggerKey, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
