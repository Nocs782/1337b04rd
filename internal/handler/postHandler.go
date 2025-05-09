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
	postService    *service.PostService
	commentService *service.CommentService
	imageStorage   domain.ImageStorage
}

func NewPostHandler(postService *service.PostService, commentService *service.CommentService, imageStorage domain.ImageStorage) *PostHandler {
	return &PostHandler{
		postService:    postService,
		commentService: commentService,
		imageStorage:   imageStorage,
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

	_, err = p.postService.CreatePost(post)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (p *PostHandler) GetActivePostsHandler(w http.ResponseWriter, r *http.Request) {
	// Call the service to get active posts
	posts, err := p.postService.GetActivePosts()
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

	post, err := p.postService.GetPostByID(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	comments, err := p.commentService.GetCommentsByPostID(post.ID)
	if err != nil {
		http.Error(w, "Failed to load comments", http.StatusInternalServerError)
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

	// Map domain.Comment to template-friendly format
	for _, c := range comments {
		data.Comments = append(data.Comments, struct {
			ID        int
			AvatarURL string
			Username  string
			Text      string
			ReplyToID *int
		}{
			ID:        c.ID,
			AvatarURL: c.AvatarURL,
			Username:  c.Author,
			Text:      c.Content,
			ReplyToID: c.ParentCommentID,
		})
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
