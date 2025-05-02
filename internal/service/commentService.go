package service

import (
	"1337b04rd/internal/adapter/postgres"
	"1337b04rd/internal/domain"
	"errors"
	"time"
)

type CommentService struct {
	repo *postgres.CommentsRepo
}

func NewCommentService(repo *postgres.CommentsRepo) *CommentService {
	return &CommentService{repo: repo}
}

// CreateComment creates a new comment
func (s *CommentService) CreateComment(comment domain.Comment) error {
	// Validate comment fields
	if comment.Content == "" || comment.Author == "" {
		return errors.New("content and author cannot be empty")
	}
	// Set creation time for the comment
	comment.CreatedAt = time.Now()

	// Call repository method to insert the comment into the database
	return s.repo.CreateComment(comment)
}

// GetCommentsByPostID fetches all comments related to a particular post
func (s *CommentService) GetCommentsByPostID(postID int) ([]domain.Comment, error) {
	// Retrieve comments from the repository
	return s.repo.GetCommentsByPostID(postID)
}
