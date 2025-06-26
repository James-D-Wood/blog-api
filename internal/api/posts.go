package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/James-D-Wood/blog-api/internal/db"
	"github.com/James-D-Wood/blog-api/internal/httputils"
	"github.com/James-D-Wood/blog-api/internal/model"
)

// TODO: implement API handlers for blog posts

func (app *App) FetchBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	post, err := app.BlogService.FetchBlogPost(r.Context(), postID)
	if err != nil {
		app.Logger.Error("failed to fetch blog post", "error", err, "location", "FetchBlogPostHandler")
		httputils.RespondWithJsonError(w, "failed to fetch blog post", 404)
		return
	}

	userID, err := httputils.GetUserFromContext(r.Context())
	if err != nil {
		app.Logger.Warn("failed to identify user", "error", err, "location", "FetchBlogPostHandler")
	}

	if post.Status == model.DRAFT {
		if post.AuthorID != userID {
			app.Logger.Error("user not authorized to view blog post", "error", err, "location", "FetchBlogPostHandler")
			httputils.RespondWithJsonError(w, "user not authorized to view this post", 403)
			return
		}
	}

	type Response struct {
		Post model.BlogPost `json:"post"`
	}

	httputils.RespondWithJson(w, Response{
		Post: post,
	})
}

func (app *App) FetchBlogPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := app.BlogService.FetchPublishedBlogPosts(r.Context())
	if err != nil {
		app.Logger.Error("failed to fetch blogs", "error", err, "location", "FetchBlogPostsHandler")
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
		app.Logger.Error("failed to read blog post payload", "error", err, "location", "CreateBlogPostHandler")
		httputils.RespondWithJsonError(w, "invalid request body", 400)
		return
	}

	userID, err := httputils.GetUserFromContext(r.Context())
	if err != nil || userID == "" {
		app.Logger.Error("failed to identify user", "error", err, "location", "CreateBlogPostHandler")
		httputils.RespondWithJsonError(w, "internal service error", 500)
		return
	}

	err = app.BlogService.CreateBlogPost(r.Context(), userID, &post)
	if err != nil {
		switch err {
		case db.ErrBlogPostAlreadyExists:
			app.Logger.Error("failed to persist blog post", "error", err, "location", "CreateBlogPostHandler")
			httputils.RespondWithJsonError(w, "cannot create blog post - resource already exists", 400)
			return
		default:
			app.Logger.Error("failed to persist blog post", "error", err, "location", "CreateBlogPostHandler")
			httputils.RespondWithJsonError(w, "internal service error", 500)
			return
		}

	}

	type Response struct {
		Post model.BlogPost `json:"post"`
	}

	httputils.RespondWithJson(w, Response{
		Post: post,
	})
}

func (app *App) UpdateBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	postID := r.PathValue("id")

	var revisedPost model.BlogPost
	err := json.NewDecoder(r.Body).Decode(&revisedPost)
	if err != nil {
		app.Logger.Error("failed to read blog post payload", "error", err, "location", "UpdateBlogPostHandler")
		httputils.RespondWithJsonError(w, "invalid request body", 400)
		return
	}

	storedPost, err := app.BlogService.FetchBlogPost(r.Context(), postID)
	if err != nil {
		app.Logger.Error("blog post for given ID does not exist", "error", err, "location", "UpdateBlogPostHandler")
		httputils.RespondWithJsonError(w, fmt.Sprintf("invalid request: blog post with ID %s does not exist", postID), 400)
		return
	}

	// validate user owns resource
	userID, err := httputils.GetUserFromContext(r.Context())
	if err != nil {
		app.Logger.Error("failed to identify user", "error", err, "location", "UpdateBlogPostHandler")
		httputils.RespondWithJsonError(w, "internal service error", 500)
		return
	}

	if storedPost.AuthorID != userID {
		app.Logger.Error("requestor does not own the blog post they are editing", "location", "UpdateBlogPostHandler", "originalAuthor", storedPost.AuthorID, "requestor", userID)
		httputils.RespondWithJsonError(w, "not authorized to update this resource", 403)
		return
	}

	err = app.BlogService.UpdateBlogPost(r.Context(), &revisedPost, &storedPost)
	if err != nil {
		// TODO: typed errors for better client responses
		app.Logger.Error("failed to persist blog post updates", "error", err, "location", "UpdateBlogPostHandler")
		httputils.RespondWithJsonError(w, "internal service error", 500)
		return
	}

	type Response struct {
		Post model.BlogPost `json:"post"`
	}

	httputils.RespondWithJson(w, Response{
		Post: storedPost,
	})
}

func (app *App) DeleteBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")

	storedPost, err := app.BlogService.FetchBlogPost(r.Context(), postID)
	if err != nil {
		app.Logger.Error("blog post for given ID does not exist", "error", err, "location", "DeleteBlogPostHandler")
		httputils.RespondWithJsonError(w, fmt.Sprintf("invalid request: blog post with ID %s does not exist", postID), 400)
		return
	}

	// validate user owns resource
	userID, err := httputils.GetUserFromContext(r.Context())
	if err != nil {
		app.Logger.Error("failed to identify user", "error", err, "location", "DeleteBlogPostHandler")
		httputils.RespondWithJsonError(w, "internal service error", 500)
		return
	}

	if storedPost.AuthorID != userID {
		app.Logger.Error("requestor does not own the blog post they are editing", "location", "DeleteBlogPostHandler", "originalAuthor", storedPost.AuthorID, "requestor", userID)
		httputils.RespondWithJsonError(w, "not authorized to update this resource", 403)
		return
	}

	app.BlogService.DeleteBlogPost(r.Context(), postID)

	type Response struct {
		PostID string `json:"post_id"`
	}

	httputils.RespondWithJson(w, Response{
		PostID: postID,
	})
}

func (app *App) AdminDeleteBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")

	_, err := app.BlogService.FetchBlogPost(r.Context(), postID)
	if err != nil {
		app.Logger.Error("blog post for given ID does not exist", "error", err, "location", "AdminDeleteBlogPostHandler")
		httputils.RespondWithJsonError(w, fmt.Sprintf("invalid request: blog post with ID %s does not exist", postID), 400)
		return
	}

	app.BlogService.DeleteBlogPost(r.Context(), postID)

	type Response struct {
		PostID string `json:"post_id"`
	}

	httputils.RespondWithJson(w, Response{
		PostID: postID,
	})
}
