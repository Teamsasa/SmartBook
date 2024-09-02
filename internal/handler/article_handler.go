package handler

import (
	"SmartBook/internal/usecase"
	"html/template"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	articleUseCase *usecase.ArticleUseCase
	template       *template.Template
}

func NewArticleHandler(articleUseCase *usecase.ArticleUseCase) *ArticleHandler {
	tmpl := template.Must(template.ParseFiles("templates/articles.html"))
	return &ArticleHandler{
		articleUseCase: articleUseCase,
		template:       tmpl,
	}
}

func (h *ArticleHandler) renderTemplate(c echo.Context, articles []usecase.Article, title string) error {
	data := struct {
		Title    string
		Articles []usecase.Article
	}{
		Title:    title,
		Articles: articles,
	}
	return h.template.Execute(c.Response().Writer, data)
}

func (h *ArticleHandler) GetLatestArticles(c echo.Context) error {
	articles, err := h.articleUseCase.GetLatestArticles()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return h.renderTemplate(c, articles, "Latest Articles")
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
	return h.renderTemplate(c, articles, "Recommended Articles")
}

// func (h *ArticleHandler) GetArticleContent(c echo.Context) error {
// 	url := c.QueryParam("url")
// 	if url == "" {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL parameter is required"})
// 	}
// 	content, err := h.articleUseCase.GetArticleContent(url)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
// 	}
// 	return c.HTML(http.StatusOK, content)
// }
