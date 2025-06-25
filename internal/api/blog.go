package api

import "net/http"

// TODO: implement API handlers for blog posts

func (app *App) FetchBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (app *App) FetchBlogPostsHandler(w http.ResponseWriter, r *http.Request) {
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
