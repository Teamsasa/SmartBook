package main

import (
	"SmartBook/internal/database"
	"SmartBook/internal/model"
	"fmt"
	"log"
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
}
