package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	// socket is the web socket for the clieent
	socket *websocket.Conn
	// recieve is a channel to recieve messages from other clients.
	recieve chan []byte
	//room is the room this client is chatting in.
	room *Room
}

func (c *Client) read() {
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *Client) Write() {
	defer c.socket.Close()
	for msg := range c.recieve {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
