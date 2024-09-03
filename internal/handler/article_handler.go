package handler

import (
	"SmartBook/internal/usecase"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	articleUseCase *usecase.ArticleUseCase
}

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
	interests := c.QueryParam("interests")
	if interests == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Interests parameter is required"})
	}
	userInterests := strings.Split(interests, ",")

	articles, err := h.articleUseCase.GetRecommendedArticles(ctx, userInterests)
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
