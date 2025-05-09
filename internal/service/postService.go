package service

import (
	"1337b04rd/internal/adapter/postgres"
	"1337b04rd/internal/domain"
	"errors"
	"time"
)

type PostService struct {
	repo *postgres.PostRepo
}
type PostServiceInterfaces interface {
	CreatePost(post domain.Post) (domain.Post, error)
	GetPostById(postId string) (domain.Post, error)
	GetActivePosts() ([]domain.Post, error)
	GetArchivePosts() ([]domain.Post, error)
	ExpireOldPosts() error
}

func NewPostService(repo *postgres.PostRepo) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(post domain.Post) (int, error) {
	if post.Title == "" || post.Content == "" {
		return 0, errors.New("title and content cannot be empty")
	}
	post.CreatedAt = time.Now()
	post.LastCommented = time.Now()

	return s.repo.CreatePost(post)
}

func (s *PostService) GetPostByID(id int) (domain.Post, error) {

	return s.repo.GetPost(id)
}

func (s *PostService) GetActivePosts() ([]domain.Post, error) {

	return s.repo.GetActivePosts()
}

func (s *PostService) GetArchivePosts() ([]domain.Post, error) {

	return s.repo.GetArchivePosts()
}
func (s *PostService) ExpireOldPosts() error {
	return s.repo.ExpireOldPosts()
}
