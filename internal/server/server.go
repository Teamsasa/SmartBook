package server

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"SmartBook/internal/database"
	"SmartBook/internal/handler"
	"SmartBook/internal/usecase"
)

type Server struct {
	port           int
	db             *gorm.DB
	articleHandler *handler.ArticleHandler
	memoHandler    *handler.MemoHandler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.NewDB()
	articleUseCase := usecase.NewArticleUseCase()
	articleHandler := handler.NewArticleHandler(articleUseCase)
	memoUseCase := usecase.NewMemoUseCase(db)
	memoHandler := handler.NewMemoHandler(memoUseCase)

	newServer := &Server{
		port:           port,
		db:             db,
		articleHandler: articleHandler,
		memoHandler:    memoHandler,
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
