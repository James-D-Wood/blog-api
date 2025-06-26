package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/James-D-Wood/blog-api/internal/constant"
	"github.com/James-D-Wood/blog-api/internal/httputils"
)

// AuthMiddleware verifies the user is valid and passes the auth claims on in the context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: this should be more durable
		logger, _ := r.Context().Value(constant.LoggerKey).(*slog.Logger)

		// hacky way to declare which paths to bypass
		excludedPaths := map[string]bool{
			"/api/v1/login": true,
		}

		ctx := r.Context()

		// verify auth token and reject request if user ID cannot be established
		if _, ok := excludedPaths[r.URL.Path]; !ok {
			token, err := httputils.DecodeBearerAuth(r)
			if err != nil {
				logger.Error("could not authenticate user", "error", err, "location", "AuthMiddleware")
				httputils.RespondWithJsonError(w, "could not authenticate user", 401)
				return
			}

			var claims httputils.AuthClaims
			err = httputils.ExtractJWTClaims(token, &claims)
			if err != nil {
				logger.Error("could not authenticate user", "error", err, "location", "AuthMiddleware")
				httputils.RespondWithJsonError(w, "could not authenticate user", 401)
				return
			}
			if claims.UserID == "" {
				logger.Error("user ID came out empty", "location", "AuthMiddleware")
				httputils.RespondWithJsonError(w, "could not authenticate user", 401)
				return
			}
			logger.Debug("claims extracted from auth JWT", "claims", claims)

			ctx = context.WithValue(ctx, constant.UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, constant.AdminKey, claims.IsAdmin)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
