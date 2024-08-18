package chatserver

import (
	"fmt"
	"log"
	"sync"

	"github.com/SurkovIlya/chat-app/internal/models"
	"github.com/SurkovIlya/chat-app/internal/room"
	"github.com/gorilla/websocket"
)

type DataStorage interface {
	SaveRoom(roomName, userName string) error
	SaveMembersChat(roomName, userName string) error
	SaveMsg(roomName, userName, msg string) error
	GetMsgs(roomName string) ([]models.RoomMsg, error)
}

type ChatServer struct {
	Rooms   map[string]*room.Room
	Storage DataStorage
	sync.RWMutex
}

func New(storage DataStorage) *ChatServer {
	return &ChatServer{
		Rooms:   make(map[string]*room.Room),
		Storage: storage,
	}
}

func (chs *ChatServer) AddRoom(roomName string, client *websocket.Conn, userName string) {
	r := room.NewRoom(client)

	chs.Lock()
	chs.Rooms[roomName] = r
	chs.Unlock()

	err := chs.Storage.SaveRoom(roomName, userName)
	if err != nil {
		log.Printf("error AddRoom: %s", err)
	}
}

func (chs *ChatServer) GetAllRooms() []string {
	roomsNames := make([]string, 0, len(chs.Rooms))

	chs.Lock()
	defer chs.Unlock()

	for name := range chs.Rooms {
		roomsNames = append(roomsNames, name)
	}

	return roomsNames
}

func (chs *ChatServer) JoinRoom(roomName string, user models.User) error {
	chs.RLock()
	defer chs.RUnlock()

	room, ok := chs.Rooms[roomName]
	if !ok {
		return fmt.Errorf("room doesn't exist")
	}

	welcome := fmt.Sprintf("!!!Пользователь %s присоединился к комнате!!!", user.UserName)

	err := chs.WriteMsg(roomName, "SERVER", welcome)
	if err != nil {
		return fmt.Errorf("error WriteMsg: %s", err)
	}

	room.JoinRoom(user.Conn)
	oldMsgs, err := chs.Storage.GetMsgs(roomName)
	if err != nil {
		return fmt.Errorf("error GetMsgs: %s", err)
	}

	for _, msg := range oldMsgs {
		user.Receive <- []byte(fmt.Sprintf("->%s<- %s: %s", roomName, msg.UserName, msg.Content))
	}

	err = chs.Storage.SaveMembersChat(roomName, user.UserName)
	if err != nil {
		return fmt.Errorf("error SaveMembersChat: %s", err)
	}

	return nil
}

func (chs *ChatServer) WriteMsg(roomName, userName, msg string) error {
	chs.RLock()
	defer chs.RUnlock()

	room, ok := chs.Rooms[roomName]
	if !ok {
		return fmt.Errorf("room doesn't exist")
	}

	room.WriteMsg(roomName, userName, msg)

	if userName != "SERVER" {
		err := chs.Storage.SaveMsg(roomName, userName, msg)
		if err != nil {
			return fmt.Errorf("error SaveMsg: %s", err)
		}
	}

	return nil
}
