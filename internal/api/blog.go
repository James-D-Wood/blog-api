package api

import (
	"encoding/json"
	"net/http"

	"github.com/James-D-Wood/blog-api/internal/httputils"
	"github.com/James-D-Wood/blog-api/internal/model"
)

// TODO: implement API handlers for blog posts

func (app *App) FetchBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	post, err := app.BlogService.FetchBlog(r.Context(), postID)
	if err != nil {
		app.Logger.Error("FetchBlogPostHandler: failed to fetch blog post", "error", err)
		httputils.RespondWithJsonError(w, "failed to fetch blog post", 404)
		return
	}

	type Response struct {
		Post model.BlogPost `json:"post"`
	}

	httputils.RespondWithJson(w, Response{
		Post: post,
	})
}

func (app *App) FetchBlogPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := app.BlogService.FetchPublishedBlogs(r.Context())
	if err != nil {
		app.Logger.Error("FetchBlogPostsHandler: failed to fetch blogs", "error", err)
		httputils.RespondWithJsonError(w, "failed to fetch blogs", 500)
		return
	}

	type Response struct {
		Posts []model.BlogPost `json:"posts"`
	}

	httputils.RespondWithJson(w, Response{
		Posts: posts,
	})
}

func (app *App) CreateBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var post model.BlogPost
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		app.Logger.Error("CreateBlogPostHandler: failed to read blog post payload", "error", err)
		httputils.RespondWithJsonError(w, "invalid request body", 400)
		return
	}

	err = app.BlogService.CreateBlog(r.Context(), &post)
	if err != nil {
		// TODO: typed errors for better client responses
		app.Logger.Error("CreateBlogPostHandler: failed to persist blog post", "error", err)
		httputils.RespondWithJsonError(w, "internal service error", 400)
		return
	}

	type Response struct {
		Post model.BlogPost `json:"post"`
	}

	httputils.RespondWithJson(w, Response{
		Post: post,
	})
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
