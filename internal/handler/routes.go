package handler

import (
	"1337b04rd/internal/adapter/postgres"
	"1337b04rd/internal/domain"
	"1337b04rd/internal/service"
	"database/sql"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, imageStorage domain.ImageStorage) {

	postRepo := postgres.NewPostRepo(db)
	postService := service.NewPostService(postRepo)
	postHandler := NewPostHandler(postService, imageStorage)

	// commentRepo := postgres.NewCommentsRepo(db)
	// commentService := service.NewCommentService(commentRepo)
	// commentHandler := NewCommentHandler(commentService)

	mux.HandleFunc("/", ShowCatalog(postService))

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	// mux.Handle("/archive", func() {}) // loads all archive
	// mux.HandleFunc("/archive/", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("/create-post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			postHandler.GetFormPostHandler(w, r)
		} else if r.Method == http.MethodPost {
			postHandler.CreatePostHandler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/post/", postHandler.GetPostByIdHandler)

	mux.HandleFunc("/archive", postHandler.GetAllPostsHandler)
}
