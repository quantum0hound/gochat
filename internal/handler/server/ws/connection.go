package ws

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	wsConn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// write writes a message with the given message type and payload.
func (c *connection) write(messageType int, payload []byte) error {
	if err := c.wsConn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		logrus.Errorf("failed to write message to connection: %s", err.Error())
		return err
	}
	return c.wsConn.WriteMessage(messageType, payload)
}
