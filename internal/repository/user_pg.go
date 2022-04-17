package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/quantum0hound/gochat/internal/models"
)

type UserProviderPostgres struct {
	db *sqlx.DB
}

func NewUserProviderPostgres(db *sqlx.DB) *UserProviderPostgres {
	return &UserProviderPostgres{db: db}
}

func (r *UserProviderPostgres) Create(user *models.User) (int, error) {
	if user == nil {
		return 0, errors.New("incorrect args")
	}
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

func (r *UserProviderPostgres) Get(username, passwordHash string) (*models.User, error) {
	if len(username) == 0 || len(passwordHash) == 0 {
		return nil, errors.New("incorrect args")
	}
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE (username=$1 AND password_hash=$2)", usersTable)
	row := r.db.QueryRow(query, username, passwordHash)
	if err := row.Scan(&user); err == sql.ErrNoRows {
		return nil, errors.New("incorrect username or password")
	}
	return &user, nil
}

func (r *UserProviderPostgres) Exists(username string) bool {
	if len(username) == 0 {
		return false
	}
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE (username=$1)", usersTable)
	err := r.db.Get(&user, query, username)
	if err != nil {
		return false
	}
	return true
}
