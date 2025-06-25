package api

import (
	"net/http"

	"github.com/James-D-Wood/blog-api/internal/httputils"
)

// TODO: implement API handlers for blog posts

func (app *App) FetchBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (app *App) FetchBlogPostsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := app.BlogService.FetchPublishedBlogs()
	if err != nil {
		app.Logger.Error("FetchBlogPostsHandler: failed to fetch blogs", "error", err)
		httputils.RespondWithJsonError(w, "failed to fetch blogs", 500)
		return
	}
	w.Write([]byte("not implemented"))
}

func (app *App) CreateBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (app *App) UpdateBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (app *App) DeleteBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (app *App) AdminDeleteBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
