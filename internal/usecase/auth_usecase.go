package usecase

import (
	"SmartBook/internal/model"
	"SmartBook/internal/repository"

	"github.com/labstack/echo/v4"
)

type IAuthUsecase interface {
	SignUp(c echo.Context, input model.InputUser) (model.User, error)
}

type authUsecase struct {
	authRepository repository.IAuthRepository
}

func NewAuthUsecase(authRepository repository.IAuthRepository) IAuthUsecase {
	return &authUsecase{
		authRepository: authRepository,
	}
}

func (u *authUsecase) SignUp(c echo.Context, input model.InputUser) (model.User, error) {
	userId, err := u.authRepository.CreateUser(c, input)
	if err != nil {
		return model.User{}, err
	}

	res, err := u.authRepository.InsertUser(c, userId, input)
	if err != nil {
		return model.User{}, err
	}

	return res, nil
}
