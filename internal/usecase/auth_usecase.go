package usecase

import (
	"SmartBook/internal/model"
	"SmartBook/internal/repository"
	"SmartBook/internal/usecase/usecaseUtils"

	"github.com/labstack/echo/v4"
)

type IAuthUsecase interface {
	SignUp(c echo.Context, input model.InputUser) (model.User, error)
	SignIn(c echo.Context, input model.InputUser) (model.User, error)
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

func (u *authUsecase) SignIn(c echo.Context, input model.InputUser) (model.User, error) {
	user, err := u.authRepository.Login(c, input)
	if err != nil {
		return model.User{}, err
	}

	err = usecaseUtils.CreateSession(c, user.ID)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
