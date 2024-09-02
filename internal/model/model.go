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
