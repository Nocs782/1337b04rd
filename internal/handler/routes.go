package handler

import (
	"database/sql"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	var chandler CommentHandler
	mux.Handle("/", &chandler) // loads all posts
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	mux.Handle("/archive", func() {}) // loads all archive
	mux.HandleFunc("/archive/", func(w http.ResponseWriter, r *http.Request) {})

	mux.Handle("/create", func() {}) // creating post
	mux.HandleFunc("/create/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// make a post
		case http.MethodGet:
			// get form to make a post
		}
	})
	mux.Handle("/post", func() {})
	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {})

	mux.Handle("/comment", func() {})
	mux.HandleFunc("/comment/", func(w http.ResponseWriter, r *http.Request) {})
}
