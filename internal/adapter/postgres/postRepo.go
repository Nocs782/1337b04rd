package postgres

import (
	"1337b04rd/internal/domain"
	"database/sql"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{}
}

func (p *PostRepo) CreatePost(post domain.Post) (int, error) {
	query := `INSERT INTO posts(title, content, avatarurl, imgsurl, author, created_at, lastcommeted, deleted ) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`

	_, err := p.db.Exec(query, post.Title, post.Content, post.AvatarURL, post.IMGsURLs, post.Author, post.CreatedAt, post.LastCommented, post.Deleted)
	if err != nil {
		return 0, err
	}
	return 1, err

}

func (p *PostRepo) GetPost(id int) (domain.Post, error) {
	var post domain.Post
	query := `SELECT * FROM posts WHERE id = $1;`
	row := p.db.QueryRow(query, id)
	err := row.Scan(post)
	if err != nil {
		return domain.Post{}, err
	}
	return post, nil
}
func (p *PostRepo) GetActivePosts() ([]domain.Post, error) {
	var posts []domain.Post
	query := `SELECT * FROM posts WHERE deleted IS FALSE`
	rows, err := p.db.Query(query)
	if err != nil {
		return posts, err

	}
	defer rows.Close()
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(post)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)

	}
	return posts, nil

}
func (p *PostRepo) GetAllPosts() ([]domain.Post, error) {
	var posts []domain.Post
	query := `SELECT * FROM posts;`
	rows, err := p.db.Query(query)
	if err != nil {
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(post)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *PostRepo) ExpirePost(id int) error {
	query := `UPDATE posts WHERE id = $1 SET deleted IS TRUE`
	_, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
