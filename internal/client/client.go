package client

import (
	"fmt"
	"log"
	"strings"

	chatserver "github.com/SurkovIlya/chat-app/internal/chat_server"
	"github.com/SurkovIlya/chat-app/internal/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserName   string
	ChatServer *chatserver.ChatServer
	Socket     *websocket.Conn
	Receive    chan []byte
}

func New(userName string, socket *websocket.Conn, chs *chatserver.ChatServer) *Client {

	return &Client{
		UserName:   userName,
		Socket:     socket,
		Receive:    make(chan []byte),
		ChatServer: chs,
	}
}

func (c *Client) Read() {
	defer c.Socket.Close()
	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			return
		}

		commandType := strings.Split(string(msg), " ")

		switch commandType[0] {
		case "/help":
			c.Receive <- []byte(helpMsg)
		case "/check_all_rooms":
			names := c.ChatServer.GetAllRooms()
			c.Receive <- []byte(fmt.Sprintf("вот комнаты: %s", strings.Join(names, ", ")))
		case "/create_room":
			err := c.ChatServer.AddRoom(commandType[1], models.User{UserName: c.UserName, Conn: c.Socket})
			if err != nil {
				log.Printf("error create_room: %s", err)
				c.Receive <- []byte(fmt.Sprint("Возникла проблема: ", err))

				continue
			}
			c.Receive <- []byte(fmt.Sprintf("Вы создали комнату %s", commandType[1]))
		case "/join_room":
			err := c.ChatServer.JoinRoom(commandType[1], models.User{UserName: c.UserName, Conn: c.Socket, Receive: c.Receive})
			if err != nil {
				log.Printf("error JoinRoom: %s", err)
				c.Receive <- []byte(fmt.Sprint("Возникла проблема: ", err))

				continue
			}
			c.Receive <- []byte(fmt.Sprintf("Вы присоединидись к комнате %s", commandType[1]))
		case "/send":
			err := c.ChatServer.WriteMsg(commandType[1], c.UserName, strings.Join(commandType[2:], " "))
			if err != nil {
				log.Printf("error WriteMsg: %s", err)
				c.Receive <- []byte(fmt.Sprintf("Возникла проблема: %s", err))
			}
		}

	}
}

func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Receive {
		err := c.Socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

func (c *Client) Welcome() {
	c.Receive <- []byte("Добро пожаловать! Введите команду /help")
}
