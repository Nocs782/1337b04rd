package postgres

import (
	"1337b04rd/internal/domain"
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	if db == nil {
		log.Fatal("NewPostRepo: received a nil db")
	}
	return &PostRepo{
		db: db,
	}
}

func (p *PostRepo) CreatePost(post domain.Post) (int, error) {
	// SQL query to insert data into posts table
	query := `
		INSERT INTO posts(title, content, avatar_url, imgs_urls, author, created_at, last_commented, deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;
	`

	// Execute the query with the necessary values, using pq.Array for the imgs_urls
	var postID int
	err := p.db.QueryRow(query, post.Title, post.Content, post.AvatarURL, pq.Array(post.IMGsURLs), post.Author, post.CreatedAt, post.LastCommented, post.Deleted).Scan(&postID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return 0, err
	}

	// Return the inserted post ID
	return postID, nil
}
func (p *PostRepo) GetPost(id int) (domain.Post, error) {
	var post domain.Post
	query := `SELECT * FROM posts WHERE id = $1;`
	row := p.db.QueryRow(query, id)
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AvatarURL, pq.Array(&post.IMGsURLs), &post.Author, &post.CreatedAt, &post.LastCommented, &post.Deleted)
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
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AvatarURL, pq.Array(&post.IMGsURLs), &post.Author, &post.CreatedAt, &post.LastCommented, &post.Deleted)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *PostRepo) GetArchivePosts() ([]domain.Post, error) {
	var posts []domain.Post
	query := `SELECT * FROM posts WHERE deleted IS TRUE`
	rows, err := p.db.Query(query)
	if err != nil {
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AvatarURL, pq.Array(&post.IMGsURLs), &post.Author, &post.CreatedAt, &post.LastCommented, &post.Deleted)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *PostRepo) ExpireOldPosts() error {
	query := `UPDATE posts SET deleted = TRUE WHERE deleted = FALSE AND last_commented < NOW() - INTERVAL '15 minutes'`
	_, err := p.db.Exec(query)
	return err
}

func (p *PostRepo) UpdatePostLastCommented(postID int, t time.Time) error {
	_, err := p.db.Exec("UPDATE posts SET last_commented = $1 WHERE id = $2", t, postID)
	return err
}
