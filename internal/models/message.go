package models

import (
	"encoding/json"
	"time"
)

type Message struct {
	Id        int64     `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	ChannelId int       `json:"channelId" db:"channel_id"`
	UserId    int       `json:"userId" db:"user_id"`
	Posted    time.Time `json:"posted" db:"posted"`
	Modified  time.Time `json:"modified" db:"modified"`
}

func (m *Message) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}

func NewMessage(data []byte) (*Message, error) {
	var m Message
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
