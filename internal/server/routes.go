package server

import (
	"SmartBook/internal/middleware/auth"
	"SmartBook/internal/middleware/cors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	authMiddleware := auth.NewAuthMiddleware(s.store)

	api := e.Group("/api", cors.SetupCORS())
	{
		// // ユーザー関連
		user := api.Group("/users")
		{
			user.POST("/signup", s.authHandler.SignUp)
			user.POST("/signin", s.authHandler.SignIn)
		}

		// 記事関連
		article := api.Group("/articles", authMiddleware.SessionMiddleware())
		{
			article.GET("/latest", s.articleHandler.GetLatestArticles)
			article.GET("/:articleId", s.articleHandler.GetArticle)
			article.GET("/recommended", s.articleHandler.GetRecommendedArticles)
			article.GET("/search", s.articleHandler.SearchArticles)
			// article.GET("/content", s.articleHandler.GetArticleContent)
		}

		// メモ関連
		memo := api.Group("/memo", authMiddleware.SessionMiddleware())
		{
			memo.POST("/", s.memoHandler.CreateMemoHandler)             // メモを作成
			memo.GET("/:articleId", s.memoHandler.GetMemoHandler)       // メモを取得
			memo.PUT("/:articleId", s.memoHandler.UpdateMemoHandler)    // メモを更新
			memo.DELETE("/:articleId", s.memoHandler.DeleteMemoHandler) // メモを削除
			memo.GET("/list", s.memoHandler.GetMemosHandler)            // メモ一覧を取得
		}
	}

	return e
}
