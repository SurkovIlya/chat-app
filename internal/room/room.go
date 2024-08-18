package room

import (
	"fmt"
	"sync"

	"github.com/SurkovIlya/chat-app/internal/models"
	"github.com/gorilla/websocket"
)

type Room struct {
	clients []models.User
	forward chan []byte
	sync.RWMutex
}

func NewRoom(client models.User) *Room {
	clients := make([]models.User, 0)
	clients = append(clients, client)

	room := Room{
		forward: make(chan []byte),
		clients: clients,
	}

	return &room
}

func (r *Room) JoinRoom(client models.User) {
	r.Lock()
	defer r.Unlock()

	r.clients = append(r.clients, client)
}

func (r *Room) WriteMsg(roomName, userName, msg string) {
	data := fmt.Sprintf("->%s<- %s: %s", roomName, userName, msg)

	for _, client := range r.clients {
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			return
		}
	}
}

func (r *Room) Contains(user models.User) bool {
	for _, u := range r.clients {
		if u.UserName == user.UserName {
			return true
		}
	}

	return false
}
