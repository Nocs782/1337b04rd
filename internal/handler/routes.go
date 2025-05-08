package handler

import (
	"1337b04rd/internal/adapter/postgres"
	rickmorty "1337b04rd/internal/adapter/rickandmorty"
	"1337b04rd/internal/domain"
	"1337b04rd/internal/service"
	"database/sql"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, imageStorage domain.ImageStorage) {

	postRepo := postgres.NewPostRepo(db)
	postService := service.NewPostService(postRepo)
	postHandler := NewPostHandler(postService, imageStorage)

	commentRepo := postgres.NewCommentsRepo(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := NewCommentHandler(commentService)

	sessionRepo := postgres.NewSessionRepo(db)
	rickmortyClient := rickmorty.NewClient("https://rickandmortyapi.com/api", &http.Client{})

	// CATALOG: Always ensure session when user loads the main page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := EnsureSession(w, r, sessionRepo, rickmortyClient)
		if err != nil {
			http.Error(w, "Failed to establish session", http.StatusInternalServerError)
			return
		}
		ShowCatalog(postService)(w, r)
	})

	// CREATE POST
	mux.HandleFunc("/create-post", func(w http.ResponseWriter, r *http.Request) {
		_, err := EnsureSession(w, r, sessionRepo, rickmortyClient)
		if err != nil {
			http.Error(w, "Failed to establish session", http.StatusInternalServerError)
			return
		}

		if r.Method == http.MethodGet {
			postHandler.GetFormPostHandler(w, r)
		} else if r.Method == http.MethodPost {
			// Later we'll also pass session info here.
			postHandler.CreatePostHandler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// POST PAGE & COMMENTING
	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		_, err := EnsureSession(w, r, sessionRepo, rickmortyClient)
		if err != nil {
			http.Error(w, "Failed to establish session", http.StatusInternalServerError)
			return
		}

		if r.Method == http.MethodPost {
			commentHandler.ServeHTTP(w, r)
		} else if r.Method == http.MethodGet {
			postHandler.GetPostByIdHandler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// ARCHIVE PAGE
	mux.HandleFunc("/archive", func(w http.ResponseWriter, r *http.Request) {
		_, err := EnsureSession(w, r, sessionRepo, rickmortyClient)
		if err != nil {
			http.Error(w, "Failed to establish session", http.StatusInternalServerError)
			return
		}
		postHandler.GetAllPostsHandler(w, r)
	})
}
