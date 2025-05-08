package handler

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/service"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostHandler struct {
	service      *service.PostService
	imageStorage domain.ImageStorage
}

func NewPostHandler(service *service.PostService, imageStorage domain.ImageStorage) *PostHandler {
	return &PostHandler{
		service:      service,
		imageStorage: imageStorage,
	}
}

func (p *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request, session *domain.Session) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Cannot parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	text := r.FormValue("text")

	if title == "" || text == "" {
		http.Error(w, "Title and text are required", http.StatusBadRequest)
		return
	}

	var imageFilename string
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imageFilename = fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)

		err = p.imageStorage.UploadImage(file, imageFilename)
		if err != nil {
			http.Error(w, "Failed to upload image: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	post := domain.Post{
		Title:     title,
		Content:   text,
		Author:    session.ID,
		AvatarURL: session.AvatarURL,
	}

	if imageFilename != "" {
		post.IMGsURLs = []string{imageFilename}
	}

	_, err = p.service.CreatePost(post)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (p *PostHandler) GetActivePostsHandler(w http.ResponseWriter, r *http.Request) {
	// Call the service to get active posts
	posts, err := p.service.GetActivePosts()
	if err != nil {
		http.Error(w, "Failed to fetch active posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (p *PostHandler) GetPostByIdHandler(w http.ResponseWriter, r *http.Request, session *domain.Session) {

	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 3 {
		http.Error(w, "Invalid post URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathSegments[2])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := p.service.GetPostByID(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	imageFilename := ""
	if len(post.IMGsURLs) > 0 {
		imageFilename = post.IMGsURLs[0]
	}

	data := struct {
		Post struct {
			ID            int
			Title         string
			Text          string
			ImageFilename string
		}
		Comments []struct {
			ID        int
			AvatarURL string
			Username  string
			Text      string
			ReplyToID *int
		}
		SessionAvatar string
		SessionID     string
	}{
		Post: struct {
			ID            int
			Title         string
			Text          string
			ImageFilename string
		}{
			ID:            post.ID,
			Title:         post.Title,
			Text:          post.Content,
			ImageFilename: imageFilename,
		},
		Comments: []struct {
			ID        int
			AvatarURL string
			Username  string
			Text      string
			ReplyToID *int
		}{},
		SessionAvatar: session.AvatarURL,
		SessionID:     session.ID,
	}

	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render post page: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *PostHandler) GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := p.service.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to fetch active posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (p *PostHandler) GetFormPostHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create-post.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
	}
}
