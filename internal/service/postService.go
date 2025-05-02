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
	GetAllPosts() ([]domain.Post, error)
	ExpirePost(int) error
}

func NewPostService(repo *postgres.PostRepo) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(post domain.Post) (int, error) {
	// Validate post before creation (if needed)
	if post.Title == "" || post.Content == "" {
		return 0, errors.New("title and content cannot be empty")
	}

	// Assign a timestamp to the created post
	post.CreatedAt = time.Now()
	post.LastCommented = time.Now()

	// Call repository method to insert the post into the database
	return s.repo.CreatePost(post)
}

func (s *PostService) GetPostByID(id int) (domain.Post, error) {
	// Retrieve a specific post by its ID
	return s.repo.GetPost(id)
}

func (s *PostService) GetActivePosts() ([]domain.Post, error) {
	// Retrieve all active (not deleted) posts
	return s.repo.GetActivePosts()
}

func (s *PostService) GetAllPosts() ([]domain.Post, error) {
	// Retrieve all posts (including deleted ones)
	return s.repo.GetAllPosts()
}

func (s *PostService) ExpirePost(id int) error {
	// Mark a post as expired (deleted)
	return s.repo.ExpirePost(id)
}
