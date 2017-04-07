// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strconv"
	"sync"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	state *State

	// Inbound messages from the clients.
	Broadcast chan interface{}

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

type State struct {
	clients       map[*Client]bool
	nextUid       uint64
	nextMessageId uint64
	mutex         *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan interface{}, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		state: &State{
			clients:       make(map[*Client]bool),
			nextUid:       0,
			nextMessageId: 0,
			mutex:         &sync.Mutex{},
		},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.removeClient(client)
		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) GetNewUser(name string) *User {
	h.state.mutex.Lock()
	ret := &User{name, strconv.FormatUint(h.state.nextUid, 10)}
	h.state.nextUid++
	h.state.mutex.Unlock()
	return ret
}

func (h *Hub) GetNextMessageId() uint64 {
	h.state.mutex.Lock()
	ret := h.state.nextMessageId
	h.state.nextMessageId++
	h.state.mutex.Unlock()
	return ret
}

func (h *Hub) GetUsers() []*User {
	h.state.mutex.Lock()
	users := make([]*User, 0)
	for client := range h.state.clients {
		if client.user != nil {
			users = append(users, client.user)
		}
	}
	h.state.mutex.Unlock()

	return users
}

func (h *Hub) registerClient(client *Client) {
	h.state.mutex.Lock()
	h.state.clients[client] = true
	h.state.mutex.Unlock()
}

func (h *Hub) removeClient(client *Client) {
	// Remove a client
	h.state.mutex.Lock()
	var user *User
	if _, ok := h.state.clients[client]; ok {
		delete(h.state.clients, client)
		user = client.user
		close(client.send)
	}
	h.state.mutex.Unlock()
	if user != nil {
		h.Broadcast <- NewUserQuit(user)
	}
}

func (h *Hub) broadcastMessage(message interface{}) {
	h.state.mutex.Lock()
	for client := range h.state.clients {
		select {
		case client.send <- message:
		default:
			// send buffer is full!
			logging.Println("Warning, client's send channel is full:", client)
		}
	}
	h.state.mutex.Unlock()
}
