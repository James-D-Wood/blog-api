package api

import (
	"log/slog"
	"net/http"

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
	apiV1.HandleFunc("POST /posts", app.CreateBlogPostHandler)
	apiV1.HandleFunc("GET /posts/{id}", app.FetchBlogPostHandler)
	apiV1.HandleFunc("GET /posts", app.FetchBlogPostsHandler)
	apiV1.HandleFunc("PUT /posts/{id}", app.UpdateBlogPostHandler)
	apiV1.HandleFunc("DELETE /posts/{id}", app.DeleteBlogPostHandler)

	// admin
	apiV1.HandleFunc("DELETE /admin/posts/{id}", app.AdminDeleteBlogPostHandler)

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
