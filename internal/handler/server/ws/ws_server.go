package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/quantum0hound/gochat/internal/models"
	"github.com/quantum0hound/gochat/internal/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserId int
type ChannelId int

type WebSocketServer struct {
	srv *service.Service

	connections map[UserId]map[*connection]bool
	channels    map[UserId][]ChannelId

	// Inbound messages from the connections.
	broadcastCh chan models.Message

	addSubscriberCh    chan subscriber
	removeSubscriberCh chan subscriber

	protocolUpgrader websocket.Upgrader
}

func NewWebSocketServer(srv *service.Service) *WebSocketServer {
	return &WebSocketServer{
		srv:                srv,
		connections:        make(map[UserId]map[*connection]bool),
		broadcastCh:        make(chan models.Message),
		addSubscriberCh:    make(chan subscriber),
		removeSubscriberCh: make(chan subscriber),
		protocolUpgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *WebSocketServer) Run() {
	for {
		select {

		//handle new client
		case subscr := <-s.addSubscriberCh:
			{
				userConnections := s.connections[subscr.userId]
				if userConnections == nil {
					userConnections = make(map[*connection]bool)
					s.connections[subscr.userId] = userConnections
				}
				s.connections[subscr.userId][subscr.connection] = true
				logrus.Debugf("%s connected as userId=%d",
					subscr.connection.socket.RemoteAddr().String(),
					subscr.userId,
				)

			}

		// handle remove client
		case subscr := <-s.removeSubscriberCh:
			{
				userConnections := s.connections[subscr.userId]
				if userConnections != nil {
					delete(userConnections, subscr.connection)
				}
				if len(userConnections) == 0 {
					delete(s.channels, subscr.userId)
				}
			}
		// handle message broadcast
		case message := <-s.broadcastCh:
			{
				userConnections := s.connections[UserId(message.UserId)]
				data, err := message.ToBytes()
				if err != nil {
					logrus.Errorf("failed to encode json: %s", err.Error())
					break
				}
				for c := range userConnections {
					logrus.Debugf("Broadcasting to %s", c.socket.RemoteAddr().String())

					select {
					case c.send <- data:
					default:
						close(c.send)
						delete(userConnections, c)
						if len(userConnections) == 0 {
							delete(s.connections, UserId(message.UserId))
						}
					}
				}
			}

		}
	}
}

type authMessage struct {
	AccessToken string `json:"accessToken"`
}

func (s *WebSocketServer) ServePeer(w http.ResponseWriter, r *http.Request) {
	socket, err := s.protocolUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	var authMsg authMessage
	//get access token and validate it
	_, data, err := socket.ReadMessage()
	err = json.Unmarshal(data, &authMsg)
	if err != nil {
		logrus.Errorf("failed to read auth response: %s", err.Error())
		return
	}

	logrus.Debugf(authMsg.AccessToken)
	userId, err := s.srv.ParseAccessToken(authMsg.AccessToken)
	if err != nil {
		logrus.Errorf("failed to parse access token: %s", err.Error())
		return
	}

	conn := &connection{
		socket: socket,
		send:   make(chan []byte, 256),
	}
	subscr := subscriber{UserId(userId), conn}
	s.addSubscriberCh <- subscr
	go subscr.handleWrite()
	go subscr.handleRead(s.removeSubscriberCh, s.broadcastCh)

}
