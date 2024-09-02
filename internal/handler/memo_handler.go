package handler

import (
	"SmartBook/internal/model"
	"SmartBook/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MemoHandler struct {
	memoUseCase *usecase.MemoUseCase
}

func NewMemoHandler(memoUseCase *usecase.MemoUseCase) *MemoHandler {
	return &MemoHandler{
		memoUseCase: memoUseCase,
	}
}

func (h *MemoHandler) GetMemosHandler(c echo.Context) error {
	// tokenなりからユーザーIDを取得

	userID := ""

	memos, err := h.memoUseCase.GetMemos(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, memos)
}

func (h *MemoHandler) UpsertMemoHandler(c echo.Context) error {
	// tokenなりからユーザーIDを取得

	userID := ""

	var req *model.MemoRequest
	err := c.Bind(&req)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	req.UserID = userID

	if req.ArticleURL == "" || req.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "article_url and content are required"})
	}

	if err := h.memoUseCase.UpsertMemo(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, req)
}

func (h *MemoHandler) GetMemoHandler(c echo.Context) error {
	// tokenなりからユーザーIDを取得

	userID := ""

	// article_urlをクエリパラメータから取得
	req := &model.MemoRequest{}
	articleURL := c.QueryParam("article_url")
	req.ArticleURL = articleURL
	req.UserID = userID

	if req.ArticleURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "article_url is required"})
	}

	memo, err := h.memoUseCase.GetMemo(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, memo)
}

func (h *MemoHandler) DeleteMemoHandler(c echo.Context) error {
	// tokenなりからユーザーIDを取得

	userID := ""

	var req *model.MemoRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	req.UserID = userID

	if req.ArticleURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "article_url is required"})
	}

	if err := h.memoUseCase.DeleteMemo(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
