package ws

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	writeWait = 10 * time.Second
)

type connection struct {
	socket *websocket.Conn
	send   chan []byte
}

func (c *connection) write(messageType int, payload []byte) error {
	if err := c.socket.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		logrus.Errorf("failed to write message to socket: %s", err.Error())
		return err
	}
	return c.socket.WriteMessage(messageType, payload)
}
