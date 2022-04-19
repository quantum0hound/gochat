package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"strconv"
)

const (
	usersTable         = "users"
	channelsTable      = "channels"
	usersChannelsTable = "users_channels"
	messagesTable      = "messages"
)

const (
	pgErrorAlreadyExists = 23505
)

type Config struct {
	Host     string
	Port     uint
	Username string
	Password string
	DbName   string
	SslMode  string
}

func NewPostgresDb(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.DbName,
		cfg.Password,
		cfg.SslMode,
	))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getErrorCode(err error) int {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return 0
	}
	code, err := strconv.Atoi(string(pqErr.Code))
	if err != nil {
		return 0
	}
	return code

}
