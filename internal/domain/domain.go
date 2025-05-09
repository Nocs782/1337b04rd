package domain

import "time"

type Post struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AvatarURL     string    `json:"avatar_url"`
	IMGsURLs      []string  `json:"img_urls"`
	Author        string    `json:"author"`
	CreatedAt     time.Time `json:"createdAt"`
	LastCommented time.Time `json:"lastCommented"`
	Deleted       bool      `json:"deleted"`
}

type Comment struct {
	ID              int       `json:"id"`
	PostID          int       `json:"post_id"`
	ParentCommentID *int      `json:"parent_comment_id"`
	AvatarURL       string    `json:"avatar_url"`
	IMGsURLs        []string  `json:"img_urls"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"createdAt"`
	Author          string    `json:"author"`
	AuthorID        string    `json:"author_id"`
}

type Session struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}
