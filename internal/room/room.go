package room

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	clients []*websocket.Conn
	forward chan []byte
	sync.RWMutex
}

func NewRoom(client *websocket.Conn) *Room {
	clients := make([]*websocket.Conn, 0)
	clients = append(clients, client)

	room := Room{
		forward: make(chan []byte),
		clients: clients,
	}

	return &room
}

func (r *Room) JoinRoom(client *websocket.Conn) {
	r.Lock()
	defer r.Unlock()

	r.clients = append(r.clients, client)
}

func (r *Room) WriteMsg(roomName, userName, msg string) {
	data := fmt.Sprintf("->%s<- %s: %s", roomName, userName, msg)

	for _, client := range r.clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			return
		}
	}
}
