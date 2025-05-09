package service

import (
	"1337b04rd/internal/adapter/postgres"
	"1337b04rd/internal/domain"
	"errors"
	"fmt"
	"time"
)

type CommentService struct {
	repo     *postgres.CommentsRepo
	postRepo *postgres.PostRepo
}

func NewCommentService(commentRepo *postgres.CommentsRepo, postRepo *postgres.PostRepo) *CommentService {
	return &CommentService{
		repo:     commentRepo,
		postRepo: postRepo,
	}
}

func (s *CommentService) CreateComment(comment domain.Comment) error {
	if comment.Content == "" || comment.AuthorID == "" {
		return errors.New("content and author cannot be empty")
	}

	comment.CreatedAt = time.Now()

	// Create the comment
	if err := s.repo.CreateComment(comment); err != nil {
		return err
	}

	// Update the post's last_commented timestamp
	if err := s.postRepo.UpdatePostLastCommented(comment.PostID, comment.CreatedAt); err != nil {
		return fmt.Errorf("failed to update post timestamp: %w", err)
	}

	return nil
}

func (s *CommentService) ReplyComment(comment domain.Comment, replyID int) error {
	if comment.Content == "" || comment.AuthorID == "" {
		return errors.New("content and author cannot be empty")
	}

	comment.CreatedAt = time.Now()
	comment.ParentCommentID = &replyID

	if err := s.repo.CreateComment(comment); err != nil {
		return err
	}

	if err := s.postRepo.UpdatePostLastCommented(comment.PostID, comment.CreatedAt); err != nil {
		return fmt.Errorf("failed to update post timestamp: %w", err)
	}

	return nil
}

func (s *CommentService) GetCommentsByPostID(postID int) ([]domain.Comment, error) {
	return s.repo.GetCommentsByPostID(postID)
}
