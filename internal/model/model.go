package model

import "time"

type Memo struct {
	ID         int
	UserID     string
	ArticleURL string
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MemoRequest struct {
	UserID     string
	ArticleURL string `json:"article_url"`
	Content    string `json:"content"`
}

type Article struct {
	ID          string    `json:"id"`
	ExternalID  string    `json:"external_id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Score       int       `json:"score"`
	Author      string    `json:"author"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	Tags        []string  `json:"tags"`
}
