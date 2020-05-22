package main

// room is the room our clients will be chatting in
type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients
	forward chan []byte
}
