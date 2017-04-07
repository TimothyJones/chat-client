package main

import (
	"encoding/json"
	"testing"
	"time"
)

var (
	user = &User{"Name", "1"}
)

type mockClock struct{}

func (mockClock) Now() time.Time { return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC) }

func TestNewMessageAdded(t *testing.T) {
	clock = &mockClock{}
	confirmJson(
		`{"type":"messageAdded","payload":{"message":"Hello!\nThis is a message.","createdAt":"Tue Nov 10 23:00:00 +0000 2009","user":{"name":"Name","userId":"1"},"id":2}}`,
		NewMessageAdded("Hello!\nThis is a message.", user, 2),
		t)
}

func TestNewJoinSuccess(t *testing.T) {
	confirmJson(`{"type":"joinRequested","payload":{"name":"Name","userId":"1"}}`, NewJoinSuccess(user), t)
}

func TestNewStartTyping(t *testing.T) {
	confirmJson(`{"type":"userStartedTyping","payload":{"name":"Name","userId":"1"}}`, NewStartTyping(user), t)
}

func TestNewStopTyping(t *testing.T) {
	confirmJson(`{"type":"userStoppedTyping","payload":{"name":"Name","userId":"1"}}`, NewStopTyping(user), t)
}

func TestNewUserQuit(t *testing.T) {
	confirmJson(`{"type":"userLeft","payload":{"name":"Name","userId":"1"}}`, NewUserQuit(user), t)
}

func TestNewUserList(t *testing.T) {
	users := []*User{{"One", "1"}, {"Two", "2"}}
	confirmJson(`{"type":"usersRequested","payload":[{"name":"One","userId":"1"},{"name":"Two","userId":"2"}]}`, NewUserList(users), t)
}

func confirmJson(expected string, message interface{}, t *testing.T) {
	b, err := json.Marshal(message)
	if err != nil {
		t.Error(err, string(b))
	}

	if string(b) != expected {
		t.Errorf("Excpected %s but was %s", expected, string(b))
	}
}
