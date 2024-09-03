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
	userID := "62236f88-4668-c711-2a61-50888a142952" // tokenなりから取得
	articleID := c.Param("articleId")

	if articleID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "articleId is required"})
	}

	req := &model.MemoRequest{
		UserID:    userID,
		ArticleID: articleID,
	}

	// bodyからcontentを取得
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "content is required"})
	}

	if err := h.memoUseCase.UpsertMemo(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "memo created"})
}

func (h *MemoHandler) GetMemoHandler(c echo.Context) error {
	userID := "" // tokenなりから取得
	articleID := c.Param("articleId")

	if articleID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "articleId is required"})
	}

	req := &model.MemoRequest{
		UserID:    userID,
		ArticleID: articleID,
	}

	memo, err := h.memoUseCase.GetMemo(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, memo)
}

func (h *MemoHandler) DeleteMemoHandler(c echo.Context) error {
	userID := "" // tokenなりから取得
	articleID := c.Param("articleId")

	if articleID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "articleId is required"})
	}

	req := &model.MemoRequest{
		UserID:    userID,
		ArticleID: articleID,
	}

	if err := h.memoUseCase.DeleteMemo(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
