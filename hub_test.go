package main

import (
	"github.com/TimothyJones/chat-client/mock_main"
	"github.com/golang/mock/gomock"
	"github.com/luci/go-render/render"
	"sort"
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
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hub := NewHub()
	c1 := &Client{send: make(chan interface{})}

	hub.registerClient(c1)

	ml := mock_main.NewMocklogFuncs(mockCtrl)
	ml.EXPECT().Println("Warning, client's send channel is full:", c1)
	logging = ml

	message := &dummyMessage{value: "Hello"}

	hub.broadcastMessage(message)
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

func TestNextUserIDUnique(t *testing.T) {
	timeout := time.After(time.Second * 2)
	hub := NewHub()
	c := make(chan uint64, 2)

	next := func() {
		c <- hub.GetNextMessageId()
	}

	go next()
	go next()

	var rec [2]bool

	for i := 0; i < 2; i++ {
		select {
		case idx := <-c:
			if idx < 0 || idx >= 2 {
				t.Fatalf("Index %d is out of expected range 0-1", idx)
			}
			rec[idx] = true
		case <-timeout:
			t.Fatal("Died waiting for a unique UID")
		}
	}

	if !rec[0] || !rec[1] {
		t.Fatal("Indexes were not unique")
	}
}

func TestNextMessageIDUnique(t *testing.T) {
	hub := NewHub()
	name := "one"

	u1 := hub.GetNewUser(name)
	if u1.Id != "0" || u1.Name != name {
		t.Fatal("User ", u1, "did not match expectations {0,", name, "} ")
	}
	u1 = hub.GetNewUser(name)
	if u1.Id != "1" || u1.Name != name {
		t.Fatal("User ", u1, "did not match expectations {1,", name, "} ")
	}
}

type ById []*User

func (a ById) Len() int           { return len(a) }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ById) Less(i, j int) bool { return a[i].Id < a[j].Id }

func TestImmutableUserList(t *testing.T) {
	hub := NewHub()
	c1 := NewClient(hub, nil)
	c2 := NewClient(hub, nil)
	hub.registerClient(c1)
	hub.registerClient(c2)

	c1.user = hub.GetNewUser("One")
	c2.user = hub.GetNewUser("Two")

	U := hub.GetUsers()

	hub.removeClient(c1)
	hub.removeClient(c2)

	if len(U) != 2 {
		t.Fatal("Wrong number of users in the user list. Expected 2 but was", len(U))
	}

	sort.Sort(ById(U))

	expectEqualButDifferent(U[0], c1.user, t)
	expectEqualButDifferent(U[1], c2.user, t)

}

func expectEqualButDifferent(a, b *User, t *testing.T) {
	if a == b {
		t.Fatal("User pointers were identical for ", a)
	}
	if a.Name != b.Name {
		t.Fatalf("User names were different (%s and %s)", a.Name, b.Name)
	}
	if a.Id != b.Id {
		t.Fatalf("User ids were different (%s and %s)", a.Id, b.Id)
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
