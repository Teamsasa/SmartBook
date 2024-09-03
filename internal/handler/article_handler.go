package handler

import (
	"SmartBook/internal/model"
	"SmartBook/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	articleUseCase *usecase.ArticleUseCase
}

// func NewArticleHandler(articleUseCase *usecase.ArticleUseCase, userUseCase *usecase.UserUseCase) *ArticleHandler {
func NewArticleHandler(articleUseCase *usecase.ArticleUseCase) *ArticleHandler {
	return &ArticleHandler{
		articleUseCase: articleUseCase,
	}
}

func (h *ArticleHandler) GetLatestArticles(c echo.Context) error {
	ctx := c.Request().Context()
	articles, err := h.articleUseCase.GetLatestArticles(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch latest articles"})
	}
	return c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) GetArticle(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("articleId")
	article, err := h.articleUseCase.GetArticleByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Article not found"})
	}
	return c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) GetRecommendedArticles(c echo.Context) error {
	ctx := c.Request().Context()

	// 認証されたユーザーのIDを取得
	// 注意: この部分は実際の認証システムに合わせて実装する必要があります
	// userID := c.Get("user_id").(string) // 例: JWTミドルウェアでセットされたユーザーID

	// データベースからユーザー情報を取得
	// user, err := h.userUseCase.GetUserByID(ctx, userID)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user information"})
	// }
	user := &model.User{
		ID:          "1",
		Interests:   []string{"Network", "Go"},
		Likes:       []string{"Go"},
		RecentViews: []string{"1", "2"},
	}

	// ユーザーの興味が設定されていない場合のエラーハンドリング
	if len(user.Interests) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User interests are not set"})
	}

	// 推奨記事を取得
	articles, err := h.articleUseCase.GetRecommendedArticles(ctx, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch recommended articles"})
	}

	return c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) GetAllArticles(c echo.Context) error {
	ctx := c.Request().Context()
	articles, err := h.articleUseCase.GetAllArticles(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch all articles"})
	}
	return c.JSON(http.StatusOK, articles)
}
