package handler

import (
	"net/http"
)

type CommentHandler struct {
	service string //service of comment
}

func newCommentHandler(service string) *CommentHandler {
	return &CommentHandler{}
}

func (c *CommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		c.postComment(w, r)
	} else {
		//method not allowed
	}
}

func (c *CommentHandler) postComment(w http.ResponseWriter, r *http.Request) {
	//to post comment service
}
