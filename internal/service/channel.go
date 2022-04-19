package service

import (
	"github.com/quantum0hound/gochat/internal/models"
	"github.com/quantum0hound/gochat/internal/repository"
)

type ChannelService struct {
	channelProvider repository.ChannelProvider
}

func NewChannelService(channelProvider repository.ChannelProvider) *ChannelService {
	return &ChannelService{
		channelProvider: channelProvider,
	}
}
func (c *ChannelService) Create(channel *models.Channel) (int, error) {
	return c.channelProvider.Create(channel)
}

func (c *ChannelService) Delete(channelId, userId int) error {
	return c.channelProvider.Delete(channelId, userId)
}

func (c *ChannelService) GetAll(userId int) ([]models.Channel, error) {
	return c.channelProvider.GetAll(userId)
}
