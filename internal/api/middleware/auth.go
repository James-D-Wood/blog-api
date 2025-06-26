package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/James-D-Wood/blog-api/internal/constant"
	"github.com/James-D-Wood/blog-api/internal/httputils"
)

// AuthProtectedMiddleware verifies the user is valid and passes the auth claims on in the context
func AuthProtectedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: this should be more durable
		logger, _ := r.Context().Value(constant.LoggerKey).(*slog.Logger)

		ctx := r.Context()

		// verify auth token and reject request if user ID cannot be established
		token, err := httputils.DecodeBearerAuth(r)
		if err != nil {
			logger.Error("could not authenticate user", "error", err, "location", "AuthProtectedMiddleware")
			httputils.RespondWithJsonError(w, "could not authenticate user", 401)
			return
		}

		var claims httputils.AuthClaims
		err = httputils.ExtractJWTClaims(token, &claims)
		if err != nil {
			logger.Error("could not authenticate user", "error", err, "location", "AuthProtectedMiddleware")
			httputils.RespondWithJsonError(w, "could not authenticate user", 401)
			return
		}
		if claims.UserID == "" {
			logger.Error("user ID came out empty", "location", "AuthProtectedMiddleware")
			httputils.RespondWithJsonError(w, "could not authenticate user", 401)
			return
		}
		logger.Debug("claims extracted from auth JWT", "claims", claims)

		ctx = context.WithValue(ctx, constant.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, constant.AdminKey, claims.IsAdmin)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthOptionalMiddleware tries to ID a user if an auth token is provided
func AuthOptionalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: this should be more durable
		logger, _ := r.Context().Value(constant.LoggerKey).(*slog.Logger)

		ctx := r.Context()

		// verify if auth token is present
		token, err := httputils.DecodeBearerAuth(r)
		if err != nil {
			logger.Debug("auth token not passed", "error", err, "location", "AuthOptionalMiddleware")
			ctx = context.WithValue(ctx, constant.UserIDKey, "")
			ctx = context.WithValue(ctx, constant.AdminKey, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		var claims httputils.AuthClaims
		err = httputils.ExtractJWTClaims(token, &claims)
		if err != nil {
			logger.Error("token passed but could not authenticate user", "error", err, "location", "AuthOptionalMiddleware")
			httputils.RespondWithJsonError(w, "could not authenticate user", 401)
			return
		}
		if claims.UserID == "" {
			logger.Error("token passed but user ID came out empty", "location", "AuthOptionalMiddleware")
			httputils.RespondWithJsonError(w, "could not authenticate user", 401)
			return
		}
		logger.Debug("claims extracted from auth JWT", "claims", claims)

		ctx = context.WithValue(ctx, constant.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, constant.AdminKey, claims.IsAdmin)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnlyMiddleware verifies the user is an admin before proceeding
func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: this should be more durable
		logger, _ := r.Context().Value(constant.LoggerKey).(*slog.Logger)

		ctx := r.Context()

		// verify auth token and reject request if user ID cannot be established
		token, err := httputils.DecodeBearerAuth(r)
		if err != nil {
			logger.Error("could not authenticate user", "error", err, "location", "AdminOnlyMiddleware")
			httputils.RespondWithJsonError(w, "could not authenticate user", 401)
			return
		}

		var claims httputils.AuthClaims
		err = httputils.ExtractJWTClaims(token, &claims)
		if err != nil {
			logger.Error("could not authenticate user", "error", err, "location", "AdminOnlyMiddleware")
			httputils.RespondWithJsonError(w, "could not authenticate user", 401)
			return
		}
		if claims.UserID == "" {
			logger.Error("user ID came out empty", "location", "AdminOnlyMiddleware")
			httputils.RespondWithJsonError(w, "could not authenticate user", 401)
			return
		}
		logger.Debug("claims extracted from auth JWT", "claims", claims)

		if !claims.IsAdmin {
			logger.Error("user attempting to access admin-only endpoint is not an admin", "location", "AdminOnlyMiddleware")
			httputils.RespondWithJsonError(w, "user is not authorized to perform this action", 403)
			return
		}

		ctx = context.WithValue(ctx, constant.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, constant.AdminKey, claims.IsAdmin)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
