package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/James-D-Wood/blog-api/internal/httputils"
)

const (
	UserIDKey ContextKey = "user_id"
	AdminKey  ContextKey = "is_admin"
)

// AuthMiddleware verifies the user is valid and passes the auth claims on in the context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: this should be more durable
		logger, _ := r.Context().Value(LoggerKey).(*slog.Logger)

		// hacky way to declare which paths to bypass
		excludedPaths := map[string]bool{
			"/api/v1/login": true,
		}

		ctx := r.Context()

		// verify auth token and reject request if user ID cannot be established
		if _, ok := excludedPaths[r.URL.Path]; !ok {
			token, err := httputils.DecodeBearerAuth(r)
			if err != nil {
				logger.Error("AuthMiddleware: could not authenticate user", "error", err)
				httputils.RespondWithJsonError(w, "could not authenticate user", 401)
				return
			}

			var claims httputils.AuthClaims
			err = httputils.ExtractJWTClaims(token, &claims)
			if err != nil {
				logger.Error("AuthMiddleware: could not authenticate user", "error", err)
				httputils.RespondWithJsonError(w, "could not authenticate user", 401)
				return
			}

			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, AdminKey, claims.IsAdmin)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
