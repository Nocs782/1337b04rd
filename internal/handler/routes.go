package handler

import (
	"1337b04rd/internal/adapter/postgres"
	"1337b04rd/internal/domain"
	"1337b04rd/internal/service"
	"database/sql"
	"net/http"
	"slices"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, imageStorage domain.ImageStorage) {

	postRepo := postgres.NewPostRepo(db)
	postService := service.NewPostService(postRepo)
	postHandler := NewPostHandler(postService)

	commentRepo := postgres.NewCommentsRepo(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := NewCommentHandler(commentService)

	mux.HandleFunc("/", ShowCatalog(postService))

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	// mux.Handle("/archive", func() {}) // loads all archive
	// mux.HandleFunc("/archive/", func(w http.ResponseWriter, r *http.Request) {})

	mux.Handle("/post", postHandler)
	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		pathSegments := strings.Split(r.URL.Path, "/")
		switch r.Method {
		case http.MethodGet:
			if slices.Contains(pathSegments, "create") {
				postHandler.GetFormPostHandler(w, r)
			} else if slices.Contains(pathSegments, "archive") {
				postHandler.GetAllPostsHandler(w, r)
			} else {
				switch len(pathSegments) {
				case 3: // get post by ID
					postHandler.GetPostByIdHandler(w, r)
				case 2: // get active posts
					postHandler.GetActivePostsHandler(w, r)
				}
			}
		case http.MethodPost:
			postHandler.CreatePostHandler(w, r)
		}
	})

	mux.Handle("/comment", commentHandler)
	mux.HandleFunc("/comment/", func(w http.ResponseWriter, r *http.Request) {
		pathSegments := strings.Split(r.URL.Path, "/")

		switch r.Method {
		case http.MethodPost:
			if len(pathSegments) == 1 {
				commentHandler.replyComment(w, r)
			}
			commentHandler.postComment(w, r)

		case http.MethodGet:
			if len(pathSegments) == 1 { // comments/{postId}/
				commentHandler.getCommentsByPostIDHandler(w, r)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Failed to read image: "+err.Error(), http.StatusBadRequest)
			return
		}

		err = imageStorage.UploadImage(file, header.Filename)
		if err != nil {
			http.Error(w, "Failed to upload image: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Image uploaded successfully as " + header.Filename))
	})

	mux.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		if filename == "" {
			http.Error(w, "Filename required", http.StatusBadRequest)
			return
		}

		data, err := imageStorage.DownloadImage(filename)
		if err != nil {
			http.Error(w, "Failed to download image: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(data)
	})

	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		if filename == "" {
			http.Error(w, "Filename required", http.StatusBadRequest)
			return
		}

		err := imageStorage.DeleteImage(filename)
		if err != nil {
			http.Error(w, "Failed to delete image: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Image deleted successfully"))
	})
}
