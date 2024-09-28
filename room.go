package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Room struct {
	// clients holds all current clients in this roomk
	clients map[*Client]bool
	// join is a channel for clients wishing to join the room.

	join chan *Client

	// leave is a channel for clients wishing to leave room
	leave chan *Client

	// forward is a channel that holds  incoming messages that should be forwarded to the other clients.

	forward chan []byte
}

func NewRoowm() *Room {
	return &Room{
		forward: make(chan []byte),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

func (r *Room) run() {
	for {
		select {
		case Client := <-r.join:
			r.clients[Client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.recieve)
		case msg := <-r.forward:
			for client := range r.clients {
				client.recieve <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: Error")
		return
	}
	client := &Client{
		socket:  socket,
		recieve: make(chan []byte, messageBufferSize),
		room:    r,
	}
	r.join <- client
	defer func() { r.leave <- client }()

	go client.Write()
	client.read()
}
