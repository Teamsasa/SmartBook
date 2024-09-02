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
