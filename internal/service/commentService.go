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

func (s *CommentService) CreateComment(comment domain.Comment) error {

	if comment.Content == "" || comment.Author == "" {
		return errors.New("content and author cannot be empty")
	}

	comment.CreatedAt = time.Now()

	return s.repo.CreateComment(comment)
}

func (s *CommentService) GetCommentsByPostID(postID int) ([]domain.Comment, error) {

	return s.repo.GetCommentsByPostID(postID)
}

func (s *CommentService) ReplyComment(comment domain.Comment, replyID int) error {
	if comment.Content == "" || comment.Author == "" {
		return errors.New("content and author cannot be empty")
	}
	comment.CreatedAt = time.Now()

	comment.ParentCommentID = &replyID

	return s.repo.CreateComment(comment)

}
