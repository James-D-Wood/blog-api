package api

import (
	"log/slog"
	"net/http"

	"github.com/James-D-Wood/blog-api/internal/api/middleware"
	"github.com/James-D-Wood/blog-api/internal/db"
)

// App wraps all global/shared state for an instance of API
type App struct {
	UserService db.UserService
	BlogService db.BlogService
	Logger      *slog.Logger
}

func (app *App) RegisterRoutes() http.Handler {

	// V1 API routes
	apiV1 := http.NewServeMux()

	// login
	apiV1.HandleFunc("POST /login", app.LoginHandler)

	// blog posts
	apiV1.Handle("POST /posts", middleware.AuthProtectedMiddleware(http.HandlerFunc(app.CreateBlogPostHandler)))
	apiV1.Handle("GET /posts/{id}", middleware.AuthOptionalMiddleware(http.HandlerFunc(app.FetchBlogPostHandler)))
	apiV1.Handle("GET /posts", middleware.AuthOptionalMiddleware(http.HandlerFunc(app.FetchBlogPostsHandler)))
	apiV1.Handle("PUT /posts/{id}", middleware.AuthProtectedMiddleware(http.HandlerFunc(app.UpdateBlogPostHandler)))
	apiV1.Handle("DELETE /posts/{id}", middleware.AuthProtectedMiddleware(http.HandlerFunc(app.DeleteBlogPostHandler)))

	// admin
	apiV1.Handle("DELETE /admin/posts/{id}", middleware.AdminOnlyMiddleware(http.HandlerFunc(app.AdminDeleteBlogPostHandler)))

	// top level mux
	m := http.NewServeMux()

	// register healthcheck, etc..
	m.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong!"))
	})

	// register nested mux
	m.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1))

	return m
}
