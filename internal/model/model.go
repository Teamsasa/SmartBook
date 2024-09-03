package model

import "time"

type MemoRequest struct {
	UserID    string
	ArticleID string
	Content   string `json:"content"`
}

type Article struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Score     int       `json:"score"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	Source    string    `json:"source"`
	Tags      []string  `json:"tags"`
}

type InputUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
