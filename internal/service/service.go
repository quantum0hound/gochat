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
	GetUser(username, password string) (*models.User, error)
	CreateUser(user *models.User) (int, error)
	GenerateAccessToken(user *models.User) (string, error)
	GenerateAccessTokenId(id int) (string, error)
	ParseAccessToken(accessToken string) (int, error)
	CreateSession(userId int, fingerprint string) (*models.Session, error)
	RefreshSession(refreshToken, fingerprint string) (*models.Session, error)
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
		Auth:    NewAuthService(repo.UserProvider, repo.SessionProvider),
		Channel: NewChannelService(repo.ChannelProvider),
	}
}
