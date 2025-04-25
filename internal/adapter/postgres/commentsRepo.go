package postgres

import (
	"1337b04rd/internal/domain"
	"database/sql"
)

type CommentsRepo struct {
	db *sql.DB
}

func (c *CommentsRepo) NewCommentsRepo(com domain.Comment) error {
	query := `
	INSERT INTO comments(post_id, parent_id, avatarurl, imgsurl, content, created_at, author)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := c.db.Exec(query, com.PostID, com.ParentCommentID, com.AvatarURL, com.IMGsURLs, com.Content, com.CreatedAt, com.Author)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommentsRepo) GetCommentsByPostID(ID int) ([]domain.Comment, error) {
	query := `SELECT * FROM comments WHERE post_id = $1`

	rows, err := c.db.Query(query, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		comments = append(comments, comment)
	}
	return comments, nil
}
