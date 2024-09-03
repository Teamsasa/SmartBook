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
	articles, err := h.articleUseCase.GetLatestArticles()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) GetArticle(c echo.Context) error {
	id := c.Param("articleId")
	article, err := h.articleUseCase.GetArticleByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Article not found"})
	}
	return c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) GetRecommendedArticles(c echo.Context) error {
	interests := c.QueryParam("interests")
	userInterests := strings.Split(interests, ",")

	articles, err := h.articleUseCase.GetRecommendedArticles(userInterests)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, articles)
}