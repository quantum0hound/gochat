package ws

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type subscription struct {
	conn      *connection
	channelId ChannelId
}

// readPump pumps messages from the websocket connection to the hub.
func (s *subscription) handleRead(unregister chan subscription, broadcast chan message) {
	conn := s.conn
	defer func() {
		unregister <- *s
		if err := conn.wsConn.Close(); err != nil {
			logrus.Errorf("Failed to close websocket connection : %s", err.Error())
		}
	}()

	conn.wsConn.SetReadLimit(maxMessageSize)
	if err := conn.wsConn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		logrus.Errorf("Failed set read deadline : %s", err.Error())
		return
	}
	conn.wsConn.SetPongHandler(
		func(string) error { err := conn.wsConn.SetReadDeadline(time.Now().Add(pongWait)); return err })
	for {
		_, msg, err := conn.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logrus.Errorf("websocket read error: %v", err)
			}
			break
		}
		m := message{msg, s.channelId}
		broadcast <- m
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) handleWrite() {
	conn := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := conn.wsConn.Close(); err != nil {
			logrus.Errorf("Failed to close websocket connection : %s", err.Error())
		}
	}()
	for {
		select {
		case message, ok := <-conn.send:
			if !ok {
				if err := conn.write(websocket.CloseMessage, []byte{}); err != nil {
					logrus.Errorf("Failed to send close message to channel: %d : %s", s.channelId, err.Error())
				}
				return
			}
			if err := conn.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := conn.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
