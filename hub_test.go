package main

import (
	"github.com/luci/go-render/render"
	"testing"
	"time"
)

func TestSmokeNewHub(t *testing.T) {
	hub := NewHub()

	if hub == nil {
		t.Fail()
	}
}

func TestRegisterClient(t *testing.T) {
	hub := NewHub()
	c := &Client{}

	hub.registerClient(c)

	if hub.state.clients[c] != true {
		t.Error("Client was not registerd")
	}
}

func TestRemoveClientNilUser(t *testing.T) {
	hub := NewHub()
	c := NewClient(hub, nil)
	hub.state.clients[c] = true

	hub.removeClient(c)
	if val, ok := hub.state.clients[c]; ok {
		t.Error("Client was not removed from map: ", val)
	}

	_, ok := (<-c.send)
	if ok {
		t.Error("Client channel is not closed")
	}

	select {
	case m := <-hub.Broadcast:
		t.Error("Unexpected broadcast from hub: ", m)
	default:
	}
}

func TestRemoveClientWithUser(t *testing.T) {
	hub := NewHub()
	c := NewClient(hub, nil)
	hub.state.clients[c] = true

	expectedName := "Test Name"
	expectedId := "1"
	c.user = &User{expectedName, expectedId}

	hub.removeClient(c)
	if val, ok := hub.state.clients[c]; ok {
		t.Error("Client was not removed from map: ", val)
	}

	_, ok := (<-c.send)
	if ok {
		t.Error("Client channel is not closed")
	}

	select {
	case m := <-hub.Broadcast:
		message, ok := m.(*userResponse)
		if !ok {
			t.Error("Wrong response type from the broadcast:", m)
		} else {
			if message.Type != USER_LEFT || message.Payload == nil || message.Payload.Name != expectedName || message.Payload.Id != expectedId {
				t.Error("Message details were not correct", render.Render(message))
			}
		}
	default:
		t.Error("No remove user broadcast from hub")
	}
}

type dummyMessage struct {
	value string
}

func TestBroadcastIntegration(t *testing.T) {
	hub := NewHub()
	c1 := NewClient(hub, nil)
	c2 := NewClient(hub, nil)
	hub.registerClient(c1)
	hub.registerClient(c2)

	message := &dummyMessage{value: "Hello"}

	hub.broadcastMessage(message)

	expectMessage(message, c1, t)
	expectMessage(message, c2, t)
}

func TestBroadcastFullClientIntegration(t *testing.T) {
	hub := NewHub()
	c1 := &Client{send: make(chan interface{})}

	hub.registerClient(c1)

	message := &dummyMessage{value: "Hello"}

	hub.broadcastMessage(message)

	expectMessage(message, c1, t)
}

type internalHub interface {
	broadcast(interface{})
	registerClient(*Client)
	removeClient(*Client)
	run()
}

func TestRunIntegration(t *testing.T) {
	hub := NewHub()
	c := NewClient(hub, nil)

	timeout := time.After(time.Second * 2)
	go hub.Run()

	select {
	case hub.Register <- c:
	case <-timeout:
		t.Fatal("Died waiting to register with the hub")
	}

	message := &dummyMessage{value: "Hello"}
	hub.Broadcast <- message

	select {
	case m := <-c.send:
		if m != message {
			t.Error("Unexpected message on client", render.Render(m))
		}
	case <-timeout:
		t.Fatal("Died waiting for client message")
	}
	select {
	case hub.Unregister <- c:
	case <-timeout:
		t.Fatal("Died waiting to unregister with the hub")
	}

	_, ok := <-c.send
	if ok {
		t.Error("Client channel was not closed")
	}
}

func expectMessage(message *dummyMessage, c *Client, t *testing.T) {
	select {
	case m := <-c.send:
		if m != message {
			t.Error("Unexpected message on client", render.Render(m))
		}
	default:
		t.Error("No broadcast recieved by client")
	}
}
