package repository

import (
	"SmartBook/internal/model"
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"firebase.google.com/go/auth"
)

type IAuthRepository interface {
	CreateUser(c echo.Context, user model.InputUser) (string, error)
	InsertUser(c echo.Context, userId string, user model.InputUser) (model.User, error)
}

type AuthRepository struct {
	db   *gorm.DB
	auth *auth.Client
}

func NewAuthRepository(db *gorm.DB, auth *auth.Client) (*AuthRepository, error) {
	return &AuthRepository{
		db:   db,
		auth: auth,
	}, nil
}

func (r *AuthRepository) CreateUser(c echo.Context, user model.InputUser) (string, error) {
	params := (&auth.UserToCreate{}).
		Email(user.Email).
		EmailVerified(false).
		Password(user.Password).
		DisplayName(user.Name)

	createdUser, err := r.auth.CreateUser(context.Background(), params)
	if err != nil {
		return "", err
	}

	return createdUser.UID, nil
}

func (r *AuthRepository) InsertUser(c echo.Context, userId string, user model.InputUser) (model.User, error) {
	newUser := model.User{
		ID:    userId,
		Name:  user.Name,
		Email: user.Email,
	}

	if err := r.db.Create(&newUser).Error; err != nil {
		return model.User{}, err
	}

	return newUser, nil
}
