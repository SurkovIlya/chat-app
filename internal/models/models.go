package models

import "github.com/gorilla/websocket"

type RoomMsg struct {
	UserName string
	Content  string
}

type User struct {
	UserName string
	Conn     *websocket.Conn
	Receive  chan []byte
}
