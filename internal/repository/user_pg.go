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
		if pgErrorAlreadyExists == getErrorCode(err) {
			err = errors.New("user already exists")
		}
		return 0, err
	}
	return id, nil
}

func (r *UserProviderPostgres) Get(username, passwordHash string) (*models.User, error) {
	if len(username) == 0 || len(passwordHash) == 0 {
		return nil, errors.New("incorrect args")
	}
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE (username=$1 AND password_hash=$2)", usersTable)

	row := r.db.QueryRow(query, username, passwordHash)
	if err := row.Scan(&user.Id, &user.Username, &user.Password); err == sql.ErrNoRows {
		return nil, errors.New("incorrect username or password")
	} else if err != nil {

		return nil, err
	}

	return &user, nil
}

func (r *UserProviderPostgres) GetById(id int) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserProviderPostgres) Exists(username string) bool {
	if len(username) == 0 {
		return false
	}
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE (username=$1)", usersTable)
	err := r.db.Get(&id, query, username)
	if err != nil {
		return false
	}
	return true
}
