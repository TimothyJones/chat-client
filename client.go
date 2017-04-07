// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan interface{}

	user *User
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{hub: hub, conn: conn, send: make(chan interface{}, 256)}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
		logging.Println("Closed reader for ", c.user)
	}()
	for {
		var f interface{}

		err := c.conn.ReadJSON(&f)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logging.Printf("error: %v", err)
			}
			break
		}
		message := f.(map[string]interface{})

		logging.Printf("Got message: %#v\n", f)
		switch message["type"] {
		case JOIN_REQUESTED:
			payload := message["payload"].(map[string]interface{})

			c.user = c.hub.GetNewUser(payload["name"].(string))
			c.hub.Broadcast <- NewUserList(c.hub.GetUsers())
			c.send <- NewJoinSuccess(c.user)
		case LIST_USERS:
			c.send <- NewUserList(c.hub.GetUsers())
		case USER_START_TYPING:
			c.hub.Broadcast <- NewStartTyping(c.user)
		case USER_STOP_TYPING:
			c.hub.Broadcast <- NewStopTyping(c.user)
		case MESSAGE_ADDED:
			body := message["payload"].(map[string]interface{})["message"].(string)
			c.hub.Broadcast <- NewMessageAdded(body, c.user, c.hub.GetNextMessageId())
		case USER_REFRESHED:
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
		logging.Println("Closed writer for ", c.user)
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			b, err := json.Marshal(message)
			if err != nil {
				logging.Println("Failed to marshal message ", message)
			}
			logging.Println("Sending message", string(b))

			err = c.conn.WriteJSON(message)
			if err != nil {
				logging.Println("Error sending message: ", err)
				return
			}
		}
	}
}
