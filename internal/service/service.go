package service

import (
	"github.com/quantum0hound/gochat/internal/models"
	"github.com/quantum0hound/gochat/internal/repository"
)

type Service struct {
	Auth
}

type Auth interface {
	CreateUser(user *models.User) (int, error)
	GenerateAccessToken(username, password string) (string, error)
	ParseAccessToken(accessToken string) (int, error)
	GenerateRefreshToken() string
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repo),
	}
}
