package room

import (
	"fmt"
	"sync"

	"github.com/SurkovIlya/chat-app/internal/models"
	"github.com/gorilla/websocket"
)

type Room struct {
	clients []models.User
	sync.RWMutex
}

const ServerName = "SERVER"

func NewRoom(client models.User) *Room {
	clients := make([]models.User, 0)

	if client.Conn != nil {
		clients = append(clients, client)
	}

	room := Room{
		clients: clients,
	}

	return &room
}

func (r *Room) JoinRoom(client models.User) {
	r.Lock()
	defer r.Unlock()

	r.clients = append(r.clients, client)
}

func (r *Room) WriteMsg(roomName, userName, msg string) error {
	data := fmt.Sprintf("->%s<- %s: %s", roomName, userName, msg)

	var userAccess bool

	r.RLock()
	userAccess = func(u string) bool {
		if u == ServerName {
			return true
		}
		for _, client := range r.clients {
			if client.UserName == u {
				return true
			}
		}

		return false
	}(userName)
	r.RUnlock()

	if !userAccess {
		return fmt.Errorf("user is not connected to the room")
	}

	for _, client := range r.clients {
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			return fmt.Errorf("error WriteMessage: %s", err)
		}
	}

	return nil
}

func (r *Room) Contains(user models.User) bool {
	for _, u := range r.clients {
		if u.UserName == user.UserName {
			return true
		}
	}

	return false
}
