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

func (c *ChannelService) SearchForChannels(pattern string) ([]models.Channel, error) {
	return c.channelProvider.SearchForChannels(pattern)
}

func (c *ChannelService) Join(channelId, userId int) (*models.Channel, error) {
	return c.channelProvider.Join(channelId, userId)
}

func (c *ChannelService) Leave(channelId, userId int) error {
	return c.channelProvider.Leave(channelId, userId)
}
