package service

import (
	"github.com/quantum0hound/gochat/internal/models"
	"github.com/quantum0hound/gochat/internal/repository"
)

type Service struct {
	Auth
	Channel
}

type Auth interface {
	CreateUser(user *models.User) (int, error)
	GenerateAccessToken(username, password string) (string, error)
	ParseAccessToken(accessToken string) (int, error)
	GenerateRefreshToken() string
}

type Channel interface {
	Create(user *models.Channel) (int, error)
	GetAll(userId int) ([]models.Channel, error)
	Delete(channelId, userId int) error
	Join(channelId, userId int) (*models.Channel, error)
	Leave(channelId, userId int) error
	SearchForChannels(pattern string) ([]models.Channel, error)
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth:    NewAuthService(repo.UserProvider),
		Channel: NewChannelService(repo.ChannelProvider),
	}
}
