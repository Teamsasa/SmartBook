package handler

import (
	"SmartBook/internal/model"
	"SmartBook/internal/usecase"
	"net/http"
	"time"

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
	userID := c.Get("userID").(string)

	memos, err := h.memoUseCase.GetMemos(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, memos)
}

func (h *MemoHandler) CreateMemoHandler(c echo.Context) error {
	userID := c.Get("userID").(string)

	type CreateMemoRequest struct {
		ArticleData *model.ArticleData `json:"article"`
		MemoContent string             `json:"content"`
	}

	var req CreateMemoRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if req.MemoContent == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "content is required"})
	}

	articleCreateReq := &model.ArticleData{
		ID:        req.ArticleData.ID,
		URL:       req.ArticleData.URL,
		Title:     req.ArticleData.Title,
		Author:    req.ArticleData.Author,
		CreatedAt: time.Now(),
	}

	memoCreateReq := &model.MemoData{
		UserID:    userID,
		ArticleID: articleCreateReq.ID,
		Content:   req.MemoContent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.memoUseCase.CreateMemo(memoCreateReq, articleCreateReq); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "memo created"})
}

func (h *MemoHandler) UpdateMemoHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
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

	if err := h.memoUseCase.UpdateMemo(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "memo updated"})
}

func (h *MemoHandler) GetMemoHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
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
	userID := c.Get("userID").(string)
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
