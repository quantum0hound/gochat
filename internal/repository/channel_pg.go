package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/quantum0hound/gochat/internal/models"
)

type ChannelProviderPostgres struct {
	db *sqlx.DB
}

func NewChannelProviderPostgres(db *sqlx.DB) *ChannelProviderPostgres {
	return &ChannelProviderPostgres{db: db}
}

// Create Adds new channel to db and associates it with user
func (r *ChannelProviderPostgres) Create(channel *models.Channel) (int, error) {
	if channel == nil {
		return 0, errors.New("unable create new channel : incorrect args")
	}
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createChannelQuery := fmt.Sprintf("INSERT INTO %s (name,description,creator) VALUES ($1, $2, $3) RETURNING id",
		channelsTable)
	row := tx.QueryRow(createChannelQuery, channel.Name, channel.Description, channel.Creator)

	if err := row.Scan(&channel.Id); err != nil {
		if pgErrorAlreadyExists == getErrorCode(err) {
			err = errors.New("channel with specified name already exists")
		}
		_ = tx.Rollback()
		return 0, err
	}

	createUsersChannelsQuery := fmt.Sprintf("INSERT INTO %s (user_id,channel_id) VALUES ($1, $2)",
		usersChannelsTable)
	_, err = tx.Exec(createUsersChannelsQuery, channel.Creator, channel.Id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return channel.Id, tx.Commit()
}

func (r *ChannelProviderPostgres) Delete(channelName string) error {
	if len(channelName) == 0 {
		return errors.New("unable to delete channel : incorrect args")
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE name = $1`, channelsTable)

	_, err := r.db.Exec(query, channelName)
	return err

}
func (r *ChannelProviderPostgres) GetById(channelId int) (*models.Channel, error) {
	var channel models.Channel
	query := fmt.Sprintf("SELECT * FROM %s WHERE (id=$1)", channelsTable)
	row := r.db.QueryRow(query, channelId)
	if err := row.Scan(&channel); err == sql.ErrNoRows {
		return nil, errors.New("channel is not found")
	}
	return &channel, nil
}
func (r *ChannelProviderPostgres) GetAll() ([]models.Channel, error) {
	return nil, errors.New("")
}

func (r *ChannelProviderPostgres) GetByName(name string) (*models.Channel, error) {
	if len(name) == 0 {
		return nil, errors.New("unable to get channel by name : incorrect args")
	}

	var channel models.Channel
	query := fmt.Sprintf("SELECT * FROM %s WHERE (name=$1)", channelsTable)
	row := r.db.QueryRow(query, name)
	if err := row.Scan(&channel); err == sql.ErrNoRows {
		return nil, errors.New("channel is not found")
	}
	return &channel, nil
}
func (r *ChannelProviderPostgres) Exists(name string) bool {
	if len(name) == 0 {
		return false
	}
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE (name=$1)", channelsTable)
	err := r.db.Get(&id, query, name)
	if err != nil {
		return false
	}
	return true
}
