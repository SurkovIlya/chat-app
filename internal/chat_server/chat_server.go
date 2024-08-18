package chatserver

import (
	"fmt"
	"log"
	"sync"

	"github.com/SurkovIlya/chat-app/internal/models"
	"github.com/SurkovIlya/chat-app/internal/room"
)

type DataStorage interface {
	SaveRoom(roomName, userName string) error
	SaveMembersChat(roomName, userName string) error
	SaveMsg(roomName, userName, msg string) error
	GetMsgs(roomName string) ([]models.RoomMsg, error)
	GetAllRooms() ([]string, error)
}

type ChatServer struct {
	Rooms   map[string]*room.Room
	Storage DataStorage
	sync.RWMutex
}

func New(storage DataStorage) *ChatServer {
	chatServer := &ChatServer{
		Rooms:   make(map[string]*room.Room),
		Storage: storage,
	}

	chatServer.prepare()

	return chatServer
}

func (chs *ChatServer) AddRoom(roomName string, client models.User) error {
	chs.Lock()

	if _, ok := chs.Rooms[roomName]; ok {
		chs.Unlock()

		return fmt.Errorf("room already exists")
	}

	r := room.NewRoom(client)

	chs.Rooms[roomName] = r
	chs.Unlock()

	err := chs.Storage.SaveRoom(roomName, client.UserName)
	if err != nil {
		log.Printf("error AddRoom: %s", err)
	}

	return nil
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

	if room.Contains(user) {
		return fmt.Errorf("user already exists")
	}

	welcome := fmt.Sprintf("!!!Пользователь %s присоединился к комнате!!!", user.UserName)

	err := chs.WriteMsg(roomName, "SERVER", welcome)
	if err != nil {
		return fmt.Errorf("error WriteMsg: %s", err)
	}

	room.JoinRoom(user)
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

func (chs *ChatServer) prepare() {
	roomsNames, err := chs.Storage.GetAllRooms()
	if err != nil {
		log.Fatalf("failed to get rooms: %s", err)
	}

	chs.Lock()
	defer chs.Unlock()

	for _, roomName := range roomsNames {
		chs.Rooms[roomName] = room.NewRoom(models.User{})
	}
}
