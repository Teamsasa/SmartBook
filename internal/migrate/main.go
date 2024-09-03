package main

import (
	"SmartBook/internal/database"
	"SmartBook/internal/model"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func main() {
	dbConn := database.NewDB()
	defer func() {
		database.CloseDB(dbConn)
		fmt.Println("🟢 Successfully migrated")
	}()

	err := dbConn.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("🔴 Error migrating User: %s", err)
	}

	err = dbConn.AutoMigrate(&model.ArticleData{})
	if err != nil {
		log.Fatalf("🔴 Error migrating ArticleData: %s", err)
	}

	err = dbConn.AutoMigrate(&model.MemoData{})
	if err != nil {
		log.Fatalf("🔴 Error migrating MemoData: %s", err)
	}

	// マイグレーション後にテストデータを挿入
	insertTestData(dbConn)
	fmt.Println("🟢 Successfully inserted test data")
}

func insertTestData(dbConn *gorm.DB) {
	// テストデータを定義
	users := []model.User{
		{ID: "0440acdf-9ff3-65ad-51fb-55e95bb230f9", Name: "Test User 1", Email: "test1@example.com"},
	}

	articles := []model.ArticleData{}

	memos := []model.MemoData{}

	// テストデータを挿入
	for _, user := range users {
		if err := dbConn.Create(&user).Error; err != nil {
			log.Fatalf("🔴 Error inserting User: %s", err)
		}
	}

	for _, article := range articles {
		if err := dbConn.Create(&article).Error; err != nil {
			log.Fatalf("🔴 Error inserting ArticleData: %s", err)
		}
	}

	for _, memo := range memos {
		if err := dbConn.Create(&memo).Error; err != nil {
			log.Fatalf("🔴 Error inserting MemoData: %s", err)
		}
	}

	fmt.Println("🟢 Successfully inserted test data")
}
