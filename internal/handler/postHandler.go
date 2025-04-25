package handler

import (
	"net/http"
	"strings"
)

type postHandler struct {
	service string //post service
}

func NewPostHandler(service string) *postHandler {
	return &postHandler{}
}

func (p *postHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")

	switch r.Method {
	case http.MethodGet:
		switch len(pathSegments) {
		case 1: // get posts
			p.GetActivePostsHandler(w, r)
		case 2: //get post/id
			p.GetPostByIdHandler(w, r)
		}
	case http.MethodPost:
		p.CreatePostHandler(w, r)
	}

}
func (p *postHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	//create post service
}
func (p *postHandler) GetActivePostsHandler(w http.ResponseWriter, r *http.Request) {
	//get posts service
}

func (p *postHandler) GetPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	//get post by id service
}

func (p *postHandler) GetArchivePostHandler(w http.ResponseWriter, r *http.Request) {
	//get archive posts
}

func (p *postHandler) GetPostFormHandler(w http.ResponseWriter, r *http.Request) {
	//get post form service
}
