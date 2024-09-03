package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api")
	{
		// // ユーザー関連
		// user := api.Group("/users")
		// {
		// 	user.GET("/users/:userId", s.getUserHandler)
		// 	user.PUT("/users/:userId", s.updateUserHandler)
		// }

		// 記事関連
		article := api.Group("/articles")
		{
			article.GET("/latest", s.articleHandler.GetLatestArticles)
			article.GET("/:articleId", s.articleHandler.GetArticle)
			article.GET("/recommended", s.articleHandler.GetRecommendedArticles)
			// article.GET("/content", s.articleHandler.GetArticleContent)
		}

		// メモ関連
		memo := api.Group("/memo")
		{
			memo.GET("/:articleId", s.memoHandler.GetMemoHandler)       // メモを取得
			memo.POST("/:articleId", s.memoHandler.UpsertMemoHandler)   // メモを作成or更新
			memo.DELETE("/:articleId", s.memoHandler.DeleteMemoHandler) // メモを削除
			memo.GET("/list", s.memoHandler.GetMemosHandler)            // メモ一覧を取得
		}
	}

	return e
}
