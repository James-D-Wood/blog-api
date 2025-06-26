package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/James-D-Wood/blog-api/internal/api"
	"github.com/James-D-Wood/blog-api/internal/api/middleware"
	"github.com/James-D-Wood/blog-api/internal/db"
)

func main() {
	// TODO: read in config using viper

	// set up logger
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	// set up app
	var blogSvc db.BlogService
	nodb := os.Getenv("NODB")
	if nodb != "" {
		blogSvc = db.NewInMemoryBlogService()
	} else {
		// set up DB connection
		panic("db not implemented")
	}

	app := api.App{
		BlogService: blogSvc,
		UserService: &db.InMemoryUserService{
			Users: db.DefaultUserMap,
		},
		Logger: logger,
	}

	// set up routing
	m := app.RegisterRoutes()

	// apply middleware
	m = middleware.LoggerMiddleware(m, app.Logger)

	// simple server setup for local testing
	server := http.Server{
		Addr:    ":8080",
		Handler: m,
	}

	logger.Info(fmt.Sprintf("listening on %s", server.Addr))
	err := server.ListenAndServe()

	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
		logger.Info("shutting down server...")
	}
}
