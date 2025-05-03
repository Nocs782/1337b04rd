package handler

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/service"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type PostHandler struct {
	service *service.PostService
}

func NewPostHandler(service *service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (p *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")

	switch r.Method {
	case http.MethodGet:
		if slices.Contains(pathSegments, "create") {
			p.GetAllPostsHandler(w, r)
		} else {
			switch len(pathSegments) {
			case 2: // get post by ID
				p.GetPostByIdHandler(w, r)
			case 1: // get active posts
				p.GetActivePostsHandler(w, r)
			}
		}
	case http.MethodPost:
		p.CreatePostHandler(w, r)
	}
}

func (p *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to create a Post
	var post domain.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service to create the post
	id, err := p.service.CreatePost(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the newly created post ID
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.Itoa(id)))
}

func (p *PostHandler) GetActivePostsHandler(w http.ResponseWriter, r *http.Request) {
	// Call the service to get active posts
	posts, err := p.service.GetActivePosts()
	if err != nil {
		http.Error(w, "Failed to fetch active posts", http.StatusInternalServerError)
		return
	}

	// Respond with the list of posts in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (p *PostHandler) GetPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the URL path
	pathSegments := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(pathSegments[1])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Call the service to get the post by ID
	post, err := p.service.GetPostByID(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Respond with the post in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (p *PostHandler) GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {

}
