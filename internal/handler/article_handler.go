package handler

import (
	"net/http"
	"strconv"
	"SmartBook/internal/usecase"

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

func (h *ArticleHandler) GetArticles(c echo.Context) error {
	articles, err := h.articleUseCase.GetArticles()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) GetArticle(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("articleId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid article ID"})
	}
	article, err := h.articleUseCase.GetArticleByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if article == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Article not found"})
	}
	return c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) GetRecommendedArticles(c echo.Context) error {
	articles, err := h.articleUseCase.GetRecommendedArticles()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) GetArticleContent(c echo.Context) error {
	url := c.QueryParam("url")
	if url == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL parameter is required"})
	}
	content, err := h.articleUseCase.GetArticleContent(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.HTML(http.StatusOK, content)
}