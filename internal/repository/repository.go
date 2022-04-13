package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/quantum0hound/gochat/internal/models"
)

type Repository struct {
	UserProvider
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserProvider: NewUserPostgres(db),
	}
}

type UserProvider interface {
	Create(user *models.User) (int, error)
	Get(username, passwordHash string) (models.User, error)
	Exists(username string) bool
}
