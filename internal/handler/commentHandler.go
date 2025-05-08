package handler

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler(service *service.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (c *CommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, session *domain.Session) {
	pathSegments := strings.Split(r.URL.Path, "/")

	switch r.Method {
	case http.MethodPost:
		if len(pathSegments) == 2 {
			c.ReplyComment(w, r, session)
		} else {
			c.PostComment(w, r, session)
		}

	case http.MethodGet:
		if len(pathSegments) == 2 { // comments/{postId}/
			c.GetCommentsByPostIDHandler(w, r)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *CommentHandler) PostComment(w http.ResponseWriter, r *http.Request, session *domain.Session) {
	var comment domain.Comment

	r.ParseForm()
	text := r.FormValue("text")
	replyTo := r.FormValue("reply_to")

	id, _ := strconv.Atoi(strings.Split(r.URL.Path, "/")[2]) // post/{id}/comment

	comment = domain.Comment{
		PostID:    id,
		Content:   text,
		CreatedAt: time.Now(),
		Author:    session.ID,
		AvatarURL: session.AvatarURL,
	}

	if replyTo != "" {
		replyID, _ := strconv.Atoi(replyTo)
		comment.ParentCommentID = &replyID
	}

	err := c.service.CreateComment(comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Comment created successfully"))
}

func (c *CommentHandler) GetCommentsByPostIDHandler(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	postID, err := strconv.Atoi(pathSegments[1])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	comments, err := c.service.GetCommentsByPostID(postID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (c *CommentHandler) ReplyComment(w http.ResponseWriter, r *http.Request, session *domain.Session) {
	var comment domain.Comment
	pathSegments := strings.Split(r.URL.Path, "/")
	parentID, err := strconv.Atoi(pathSegments[1])
	if err != nil {
		http.Error(w, "Invalid parent ID", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	comment.Author = session.Name
	comment.AvatarURL = session.AvatarURL
	comment.CreatedAt = time.Now()

	err = c.service.ReplyComment(comment, parentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Comment created successfully"))
}
