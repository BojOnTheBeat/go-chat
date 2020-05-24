package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bojonthebeat/go-trace"
	"github.com/gorilla/websocket"
)

// room is the room our clients will be chatting in
type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients
	forward chan []byte

	// join is a channel for clients wishing to join the room
	join chan *client

	// leave is a channel for clients wishing to leave the room
	leave chan *client

	// client is a map that holds all the clients currently in the room
	clients map[*client]bool

	// we have two (join and leave) channels to allow us safely add and remove clients
	// from the clients map. We don't want two goroutines to allow us modify
	// the map at the same time

	// tracer will receive trace information of activity in
	// the room
	tracer trace.Tracer
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
			r.tracer.Trace("New Client joined")

		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
			r.tracer.Trace(fmt.Sprintf("%d clients remaining", len(r.clients)))

		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", string(msg))
			// forward message to all clients
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace(" -- sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// In order to use web sockets, we must upgrade the HTTP connection using the Upgrader type,
// which is reusable so we only create one.
// THen when a request comes in via the serveHTTP method, we get the socket by calling
// upgrade.Upgrade below
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// this creates a web socket, initializes a client and has that client join the room.
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write() // run in a different thread/goroutine
	client.read()

}

// newRoom makes a new room.
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}
