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
			if len(names) == 0 {
				c.Receive <- []byte("There are no rooms yet! You can create them.")

				continue
			}

			c.Receive <- []byte(fmt.Sprintf("Available rooms: %s", strings.Join(names, ", ")))
		case "/create_room":
			err := c.ChatServer.AddRoom(commandType[1], models.User{
				UserName: c.UserName,
				Conn:     c.Socket,
			})
			if err != nil {
				log.Printf("Error create room: %s", err)
				c.Receive <- []byte(fmt.Sprint("Error create room: ", err))

				continue
			}

			c.Receive <- []byte(fmt.Sprintf("You have created a room %s", commandType[1]))
		case "/join_room":
			err := c.ChatServer.JoinRoom(commandType[1], models.User{
				UserName: c.UserName,
				Conn:     c.Socket,
				Receive:  c.Receive,
			})
			if err != nil {
				log.Printf("error JoinRoom: %s", err)
				c.Receive <- []byte(fmt.Sprint("Error join room: ", err))

				continue
			}

			c.Receive <- []byte(fmt.Sprintf("You have joined the room %s", commandType[1]))
		case "/check_my_room":
			userRooms := c.ChatServer.GetUserRooms(models.User{
				UserName: c.UserName,
				Conn:     c.Socket,
			})
			if len(userRooms) == 0 {
				c.Receive <- []byte("You are not in the same room")

				continue
			}

			c.Receive <- []byte(fmt.Sprintf("You are a member of the rooms: %s", strings.Join(userRooms, ", ")))
		case "/send":
			err := c.ChatServer.WriteMsg(commandType[1], c.UserName, strings.Join(commandType[2:], " "))
			if err != nil {
				log.Printf("error WriteMsg: %s", err)
				c.Receive <- []byte(fmt.Sprintf("Error send message: %s", err))
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
	c.Receive <- []byte("Welcome! Enter the command /help")
}
