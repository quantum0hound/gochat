package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/quantum0hound/gochat/internal/models"
)

type Repository struct {
	UserProvider
	ChannelProvider
	SessionProvider
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserProvider:    NewUserProviderPostgres(db),
		ChannelProvider: NewChannelProviderPostgres(db),
		SessionProvider: NewSessionProviderPostgres(db),
	}
}

type UserProvider interface {
	Create(user *models.User) (int, error)
	Get(username, passwordHash string) (*models.User, error)
	GetById(id int) (*models.User, error)
	Exists(username string) bool
}

type SessionProvider interface {
	Create(session *models.Session) error
	Delete(id string) error
	Get(id string) (*models.Session, error)
}

type ChannelProvider interface {
	Create(channel *models.Channel) (int, error)
	Delete(channelId, userId int) error
	Join(channelId, userId int) (*models.Channel, error)
	Leave(channelId, userId int) error
	GetAll(userId int) ([]models.Channel, error)
	SearchForChannels(pattern string) ([]models.Channel, error)
	GetById(channelId int) (*models.Channel, error)
	GetByName(name string) (*models.Channel, error)
}
