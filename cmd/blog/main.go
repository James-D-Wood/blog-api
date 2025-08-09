package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/James-D-Wood/blog-api/internal/api"
	"github.com/James-D-Wood/blog-api/internal/api/middleware"
	"github.com/James-D-Wood/blog-api/internal/config"
	"github.com/James-D-Wood/blog-api/internal/db"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// set up logger
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.Logger.GetSlogLevel()}),
	)

	// set up app
	var blogSvc db.BlogService
	if !cfg.DB.Enabled {
		logger.Info("using in-memory database")
		blogSvc = db.NewInMemoryBlogService()
	} else {
		// set up DB connection
		logger.Error("database not implemented yet")
		return fmt.Errorf("database not implemented")
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

	// server setup
	server := http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: m,
	}

	logger.Info(fmt.Sprintf("listening on %s", server.Addr))
	err = server.ListenAndServe()

	if err != nil {
		if err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		logger.Info("shutting down server...")
	}

	return nil
}
