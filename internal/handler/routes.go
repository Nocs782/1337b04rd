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
	postHandler := NewPostHandler(postService)

	commentRepo := postgres.NewCommentsRepo(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := NewCommentHandler(commentService)

	mux.HandleFunc("/", ShowCatalog(postService))

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	// mux.Handle("/archive", func() {}) // loads all archive
	// mux.HandleFunc("/archive/", func(w http.ResponseWriter, r *http.Request) {})

	mux.Handle("/post", postHandler)
	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {})

	mux.Handle("/comment", commentHandler)
	mux.HandleFunc("/comment/", func(w http.ResponseWriter, r *http.Request) {})

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
