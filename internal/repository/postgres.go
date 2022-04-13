package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	usersTable         = "users"
	channelsTable      = "channels"
	usersChannelsTable = "users_channels"
	messagesTable      = "todo_items"
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
