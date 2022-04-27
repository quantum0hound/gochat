package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/quantum0hound/gochat/internal/models"
	"github.com/sirupsen/logrus"
	"time"
)

const (

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type subscriber struct {
	userId     UserId
	connection *connection
}

// readPump pumps messages from the websocket connection to the hub.
func (s *subscriber) handleRead(removeSubscriberCh chan subscriber, broadcast chan models.Message) {
	socket := s.connection.socket
	defer func() {
		removeSubscriberCh <- *s
		if err := socket.Close(); err != nil {
			logrus.Errorf("Failed to close websocket connection : %s", err.Error())
		}
	}()

	socket.SetReadLimit(maxMessageSize)
	if err := socket.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		logrus.Errorf("Failed set read deadline : %s", err.Error())
		return
	}
	socket.SetPongHandler(
		func(string) error { err := socket.SetReadDeadline(time.Now().Add(pongWait)); return err })
	for {
		_, data, err := socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logrus.Errorf("websocket read error: %v", err)
			}
			break
		}
		var message models.Message
		err = json.Unmarshal(data, &message)
		if err != nil {
			logrus.Errorf("incorrect json supplied : %s", err.Error())
		}

		logrus.Debugf("%s", message.Content)
		broadcast <- message
	}
}

func (s *subscriber) handleWrite() {
	socket := s.connection.socket
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := socket.Close(); err != nil {
			logrus.Errorf("Failed to close websocket connection : %s", err.Error())
		}
	}()
	for {
		select {
		case message, ok := <-s.connection.send:
			if !ok {
				if err := s.connection.write(websocket.CloseMessage, []byte{}); err != nil {
					logrus.Errorf("Failed to send close message to client : %s", err.Error())
				}
				return
			}
			if err := s.connection.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := s.connection.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
