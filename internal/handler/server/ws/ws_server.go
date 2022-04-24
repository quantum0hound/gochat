package ws

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ChannelId int

type message struct {
	data      []byte
	channelId ChannelId
}

type WebSocketServer struct {
	// Registered connections.
	channels map[ChannelId]map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription

	protocolUpgrader websocket.Upgrader
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		channels:   make(map[ChannelId]map[*connection]bool),
		broadcast:  make(chan message),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		protocolUpgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (ws *WebSocketServer) Run() {
	for {
		select {
		case subscr := <-ws.register:
			connections := ws.channels[subscr.channelId]
			if connections == nil {
				connections = make(map[*connection]bool)
				ws.channels[subscr.channelId] = connections
			}
			logrus.Debugf("%s connected to channel%d", subscr.conn.wsConn.RemoteAddr().String(), subscr.channelId)
			ws.channels[subscr.channelId][subscr.conn] = true
		case subscr := <-ws.unregister:
			connections := ws.channels[subscr.channelId]
			if connections != nil {
				if _, ok := connections[subscr.conn]; ok {
					delete(connections, subscr.conn)
					close(subscr.conn.send)
					if len(connections) == 0 {
						delete(ws.channels, subscr.channelId)
					}
				}
			}
		case m := <-ws.broadcast:
			connections := ws.channels[m.channelId]
			for c := range connections {
				logrus.Debugf("Broadcasting to %s", c.wsConn.RemoteAddr().String())

				select {
				case c.send <- m.data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(ws.channels, m.channelId)
					}
				}
			}
		}
	}
}

func (ws *WebSocketServer) ServePeer(w http.ResponseWriter, r *http.Request, channelId ChannelId) {
	wsConn, err := ws.protocolUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}
	c := &connection{send: make(chan []byte, 256), wsConn: wsConn}
	s := subscription{c, channelId}
	ws.register <- s
	go s.handleRead(ws.unregister, ws.broadcast)
	go s.handleWrite()
}
