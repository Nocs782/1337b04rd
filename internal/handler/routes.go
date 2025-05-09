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

	commentRepo := postgres.NewCommentsRepo(db)
	postRepo := postgres.NewPostRepo(db)
	commentService := service.NewCommentService(commentRepo, postRepo)
	commentHandler := NewCommentHandler(commentService)
	postService := service.NewPostService(postRepo)
	postHandler := NewPostHandler(postService, commentService, imageStorage)

	sessionRepo := postgres.NewSessionRepo(db)
	rickmortyClient := rickmorty.NewClient("https://rickandmortyapi.com/api", &http.Client{})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, err := EnsureSession(w, r, sessionRepo, rickmortyClient)
		if err != nil {
			http.Error(w, "Failed to establish session", http.StatusInternalServerError)
			return
		}

		ShowCatalog(postService, session)(w, r)
	})

	mux.HandleFunc("/create-post", func(w http.ResponseWriter, r *http.Request) {
		session, err := EnsureSession(w, r, sessionRepo, rickmortyClient)
		if err != nil {
			http.Error(w, "Failed to establish session", http.StatusInternalServerError)
			return
		}
		if r.Method == http.MethodGet {
			postHandler.GetFormPostHandler(w, r)
		} else if r.Method == http.MethodPost {
			postHandler.CreatePostHandler(w, r, session)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		session, err := EnsureSession(w, r, sessionRepo, rickmortyClient)
		if err != nil {
			http.Error(w, "Failed to establish session", http.StatusInternalServerError)
			return
		}

		if r.Method == http.MethodPost {
			commentHandler.ServeHTTP(w, r, session)
		} else if r.Method == http.MethodGet {
			postHandler.GetPostByIdHandler(w, r, session)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/archive", func(w http.ResponseWriter, r *http.Request) {
		ShowArchive(postService)(w, r)
	})
}
