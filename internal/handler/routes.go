package handler

import (
	"1337b04rd/internal/adapter/postgres"
	"1337b04rd/internal/service"
	"database/sql"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	postRepo := postgres.NewPostRepo(db)
	postService := service.NewPostService(postRepo)
	postHandler := NewPostHandler(postService)

	commentRepo := postgres.NewCommentsRepo(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := NewCommentHandler(commentService)

	// mux.Handle("/", func() {}) // loads all posts
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	// mux.Handle("/archive", func() {}) // loads all archive
	// mux.HandleFunc("/archive/", func(w http.ResponseWriter, r *http.Request) {})

	mux.Handle("/post", postHandler)
	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {})

	mux.Handle("/comment", commentHandler)
	mux.HandleFunc("/comment/", func(w http.ResponseWriter, r *http.Request) {})
}
