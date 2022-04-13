package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/quantum0hound/gochat/internal/models"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	query := fmt.Sprintf(
		"INSERT INTO %s (username,password_hash) values ($1, $2) RETURNING id",
		usersTable,
	)
	row := r.db.QueryRow(
		query,
		user.Username,
		user.Password, //password hash is here
	)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, passwordHash string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE (username=$1 AND password_hash=$2)", usersTable)
	err := r.db.Get(&user, query, username, passwordHash)
	return user, err
}
