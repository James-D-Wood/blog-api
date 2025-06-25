package middleware

// TODO: implement auth middleware

// type ContextKey string

// const (
// 	LoggerKey ContextKey = "logger"
// )

// // AuthMiddleware rejects requests for users that are not properly authenticated
// func AuthMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		logger.Info(fmt.Sprintf("%s %s", r.Method, r.URL))
// 		ctx := context.WithValue(r.Context(), LoggerKey, logger)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
