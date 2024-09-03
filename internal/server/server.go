package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"

	"SmartBook/internal/cache"
	"SmartBook/internal/database"
	"SmartBook/internal/handler"
	"SmartBook/internal/usecase"
)

type Server struct {
	port           int
	db             *gorm.DB
	articleHandler *handler.ArticleHandler
	memoHandler    *handler.MemoHandler
	cache          cache.Cache
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.NewDB()

	// HTTP クライアントを作成
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// キャッシュインスタンスを作成
	cacheInstance := cache.NewInMemoryCache()
	// 1時間ごとに期限切れのアイテムを削除
	go cacheInstance.StartCleanup(1 * time.Hour)

	articleUseCase, _ := usecase.NewArticleUseCase(httpClient, cacheInstance)
	articleHandler := handler.NewArticleHandler(articleUseCase)
	memoUseCase := usecase.NewMemoUseCase(db)
	memoHandler := handler.NewMemoHandler(memoUseCase)

	newServer := &Server{
		port:           port,
		db:             db,
		articleHandler: articleHandler,
		memoHandler:    memoHandler,
		cache:          cacheInstance,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
