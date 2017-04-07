package main

import (
	"time"
)

const (
	USER_JOINED       = "userJoined"
	USER_LEFT         = "userLeft"
	JOIN_REQUESTED    = "joinRequested"
	LIST_USERS        = "usersRequested"
	USER_START_TYPING = "userStartedTyping"
	USER_STOP_TYPING  = "userStoppedTyping"
	MESSAGE_ADDED     = "messageAdded"
	USER_REFRESHED    = "userRefreshed"

	RFC2822 = "Mon Jan 02 15:04:05 -0700 2006"
)

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

var clock Clock = &realClock{}

type User struct {
	Name string `json:"name"`
	Id   string `json:"userId"`
}

type userListResponse struct {
	Type    string  `json:"type"`
	Payload []*User `json:"payload"`
}

type userResponse struct {
	Type    string `json:"type"`
	Payload *User  `json:"payload"`
}

type messageAddedResponse struct {
	Type    string `json:"type"`
	Payload struct {
		Message   string `json:"message"`
		CreatedAt string `json:"createdAt"`
		User      *User  `json:"user"`
		Id        uint64 `json:"id"`
	} `json:"payload"`
}

func NewMessageAdded(data string, user *User, id uint64) *messageAddedResponse {
	ret := &messageAddedResponse{Type: MESSAGE_ADDED}
	ret.Payload.Message = data
	ret.Payload.User = user
	ret.Payload.CreatedAt = clock.Now().Format(RFC2822)
	ret.Payload.Id = id

	return ret
}

func NewJoinSuccess(user *User) *userResponse {
	return &userResponse{Type: JOIN_REQUESTED, Payload: user}
}

func NewStartTyping(user *User) *userResponse {
	return &userResponse{Type: USER_START_TYPING, Payload: user}
}

func NewStopTyping(user *User) *userResponse {
	return &userResponse{Type: USER_STOP_TYPING, Payload: user}
}

func NewUserQuit(user *User) *userResponse {
	return &userResponse{Type: USER_LEFT, Payload: user}
}

func NewUserList(users []*User) *userListResponse {
	return &userListResponse{Type: LIST_USERS, Payload: users}
}
