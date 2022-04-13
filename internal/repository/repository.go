package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/quantum0hound/gochat/internal/models"
)

type Repository struct {
	Auth
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth: NewAuthPostgres(db),
	}
}

type Auth interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, passwordHash string) (models.User, error)
}
