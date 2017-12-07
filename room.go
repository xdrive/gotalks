package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/stretchr/objx"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

type room struct {
	// incoming messages which are to be forwarded to all the clients
	forward chan *message

	clients sync.Map
}

func newRoom() *room {

	return &room{
		forward: make(chan *message),
	}
}

func (r *room) run() {
	for msg := range r.forward {
		r.clients.Range(func(k, v interface{}) bool {
			k.(*client).send <- msg
			return true
		})
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.Join(client)
	defer func() { r.Leave(client) }()
	go client.write()
	client.read()
}

func (r *room) Join(c *client) {
	r.clients.Store(c, true)
}

func (r *room) Leave(c *client) {
	r.clients.Delete(c)
}
