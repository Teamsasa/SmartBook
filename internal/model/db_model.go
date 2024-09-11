package model

import (
	"time"
)

type User struct {
	ID          string     `json:"id" gorm:"type:varchar(255);primaryKey"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Email       string     `json:"email" gorm:"type:varchar(255);unique;not null"`
	Password    string     `json:"password" gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"not null"`
	Memos       []MemoData `json:"memos" gorm:"foreignKey:UserID"`
	Interests   []string   `json:"interests" gorm:"type:varchar(255)[]"`
	RecentViews []string   `json:"recent_views" gorm:"type:varchar(255)[]"`
	Likes       []string   `json:"likes" gorm:"type:varchar(255)[]"`
}

type ArticleData struct {
	ID        string     `json:"id" gorm:"type:varchar(255);primaryKey"`
	URL       string     `json:"url" gorm:"type:varchar(1000);not null"`
	Title     string     `json:"title" gorm:"type:varchar(255);not null"`
	Author    string     `json:"author" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null"`
	Memos     []MemoData `json:"memos" gorm:"foreignKey:ArticleID"`
}

type MemoData struct {
	ID        int         `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    string      `json:"user_id" gorm:"type:varchar(255);not null"`
	ArticleID string      `json:"article_id" gorm:"type:varchar(255);not null"`
	Content   string      `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time   `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"not null"`
	User      User        `gorm:"foreignKey:UserID"`
	Article   ArticleData `gorm:"foreignKey:ArticleID"`
}
