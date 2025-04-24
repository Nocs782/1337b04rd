package domain

import "time"

type Post struct {
	ID            int
	Title         string
	Content       string
	AvatarURL     string
	IMGsURLs      []string
	Author        string
	CreatedAt     time.Time
	LastCommented time.Time
	Deleted       bool
}

type Comment struct {
	ID              int
	PostID          int
	ParentCommentID *int
	AvatarURL       string
	IMGsURLs        []string
	Content         string
	CreatedAt       time.Time
	Author          string
}

type Session struct {
	Name      string
	IP        string
	AvatarURL string
	CreatedAt time.Time
	ExpiresAt time.Time
}
