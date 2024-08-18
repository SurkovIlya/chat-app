package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/SurkovIlya/chat-app/internal/client"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (s *Server) Connect(w http.ResponseWriter, req *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	userName := s.genUserName()
	exist, err := s.UserStorage.UserExist(userName)
	if err != nil {
		log.Printf("error UserExist: %s", err)
	}

	if exist {
		userName = "the_best_" + userName
	}

	c := client.New(userName, socket, s.ChatServer)

	err = s.UserStorage.SaveUser(userName)
	if err != nil {
		log.Printf("error SaveUser: %s", err)
	}

	go c.Write()

	c.Welcome()
	c.Read()
}

func (s *Server) genUserName() string {
	arrName := []string{"cat", "dog", "bird", "snake", "frog", "cow"}

	e := rand.Intn(len(arrName))

	return fmt.Sprint(arrName[e], time.Now().UnixNano())
}
