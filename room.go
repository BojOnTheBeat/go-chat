package main

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
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true

		case client := <-r.leave:
			//leaving
			delete(r.clients, client)
			close(client.send)

		case msg := <-r.forward:
			// forward message to all clients
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}
