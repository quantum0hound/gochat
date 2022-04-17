package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/quantum0hound/gochat/internal/models"
)

type Repository struct {
	UserProvider
	ChannelProvider
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserProvider:    NewUserProviderPostgres(db),
		ChannelProvider: NewChannelProviderPostgres(db),
	}
}

type UserProvider interface {
	Create(user *models.User) (int, error)
	Get(username, passwordHash string) (*models.User, error)
	Exists(username string) bool
}

type ChannelProvider interface {
	Create(channel *models.Channel) (int, error)
	Delete(name string) error
	GetById(channelId int) (*models.Channel, error)
	GetByName(name string) (*models.Channel, error)
	Exists(name string) bool
}
