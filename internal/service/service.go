package service

import (
	"github.com/quantum0hound/gochat/internal/models"
	"github.com/quantum0hound/gochat/internal/repository"
)

type Service struct {
	repo *repository.Repository
	Auth
}

type Auth interface {
	CreateUser(user models.User) (int, error)
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}
