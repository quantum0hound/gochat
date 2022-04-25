package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/quantum0hound/gochat/internal/models"
	"strings"
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

func (r *ChannelProviderPostgres) Delete(channelId, userId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 AND creator = $2`, channelsTable)
	_, err := r.db.Exec(query, channelId, userId)
	return err

}
func (r *ChannelProviderPostgres) GetById(channelId int) (*models.Channel, error) {
	var channel models.Channel
	query := fmt.Sprintf("SELECT * FROM %s WHERE (id=$1)", channelsTable)
	err := r.db.Get(&channel, query, channelId)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *ChannelProviderPostgres) GetByName(name string) (*models.Channel, error) {
	if len(name) == 0 {
		return nil, errors.New("unable to get channel by name : incorrect args")
	}

	var channel models.Channel
	query := fmt.Sprintf("SELECT * FROM %s WHERE (name=$1)", channelsTable)
	err := r.db.Get(&channel, query, name)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *ChannelProviderPostgres) GetAll(userId int) ([]models.Channel, error) {
	var channels []models.Channel
	query := fmt.Sprintf(
		`SELECT ct.id, ct.name, ct.creator, ct.description FROM %s AS ct 
				INNER JOIN %s AS uc on ct.id = uc.channel_id WHERE uc.user_id = $1`,
		channelsTable,
		usersChannelsTable,
	)
	err := r.db.Select(&channels, query, userId)
	return channels, err
}

func (r *ChannelProviderPostgres) SearchForChannels(pattern string) ([]models.Channel, error) {
	var channels []models.Channel
	query := fmt.Sprintf(
		`SELECT id, name, creator, description FROM %s 
				WHERE (lower(name) LIKE '%s%%')`,
		channelsTable, strings.ToLower(pattern),
	)
	err := r.db.Select(&channels, query)
	return channels, err
}

func (r *ChannelProviderPostgres) Join(channelId, userId int) (*models.Channel, error) {
	channel, err := r.GetById(channelId)
	if err != nil {
		return nil, errors.New("channel not exists")
	}
	query := fmt.Sprintf("INSERT INTO %s (user_id,channel_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		usersChannelsTable)
	_, err = r.db.Exec(query, userId, channelId)

	if err != nil {
		return nil, err
	}
	return channel, nil

}

func (r *ChannelProviderPostgres) Leave(channelId, userId int) error {

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1 AND channel_id=$2",
		usersChannelsTable)
	_, err := r.db.Exec(query, userId, channelId)

	return err

}
