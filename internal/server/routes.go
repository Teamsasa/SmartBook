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
			article.GET("", s.articleHandler.GetArticles)
			article.GET("/:articleId", s.articleHandler.GetArticle)
			article.GET("/recommended", s.articleHandler.GetRecommendedArticles)
			article.GET("/content", s.articleHandler.GetArticleContent)
		}

		// メモ関連
		memo := api.Group("/memos")
		{
			memo.GET("", s.memoHandler.GetMemosHandler)
			memo.POST("", s.memoHandler.CreateMemoHandler) // 新規作成と更新を兼ねる
			memo.GET("/:memoId", s.memoHandler.GetMemoHandler)
			memo.DELETE("/:memoId", s.memoHandler.DeleteMemoHandler)
		}
	}

	return e
}

// func (s *Server) HelloWorldHandler(c echo.Context) error {
// 	resp := map[string]string{
// 		"message": "Hello World",
// 	}

// 	return c.JSON(http.StatusOK, resp)
// }

// func (s *Server) healthHandler(c echo.Context) error {
// 	return c.JSON(http.StatusOK, s.db.Health())
// }
