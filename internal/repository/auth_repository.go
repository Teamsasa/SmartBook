package repository

import (
	"SmartBook/internal/model"
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"firebase.google.com/go/auth"
)

type IAuthRepository interface {
	CreateUser(c echo.Context, user model.InputUser) (string, error)
	InsertUser(c echo.Context, userId string, user model.InputUser) (model.User, error)
	Login(c echo.Context, user model.InputUser) (model.User, error)
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
		ID:       userId,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	if err := r.db.Create(&newUser).Error; err != nil {
		return model.User{}, err
	}

	return newUser, nil
}

func (r *AuthRepository) Login(c echo.Context, user model.InputUser) (model.User, error) {
	var userData model.User
	if err := r.db.Where("email = ?", user.Email).First(&userData).Error; err != nil {
		return model.User{}, err
	}

	if user.Password != userData.Password {
		return model.User{}, errors.New("password is incorrect")
	}

	userRecord, err := r.auth.GetUserByEmail(context.Background(), user.Email)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:    userRecord.UID,
		Email: userRecord.Email,
		Name:  userRecord.DisplayName,
	}, nil
}
