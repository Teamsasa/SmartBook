package handler

import (
	"SmartBook/internal/model"
	"SmartBook/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IAuthHandler interface {
	SignUp(c echo.Context) error
	SignIn(c echo.Context) error
}

type AuthHandler struct {
	authUseCase usecase.IAuthUsecase
}

func NewAuthHandler(authUseCase usecase.IAuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	var input model.InputUser
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	res, err := h.authUseCase.SignUp(c, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) SignIn(c echo.Context) error {
	var input model.InputUser
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	res, err := h.authUseCase.SignIn(c, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, res)
}
