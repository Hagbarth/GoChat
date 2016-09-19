package chat

import (
	"github.com/satori/go.uuid"
	"time"
)

type MessageRequest struct {
	Uid     string
	Message string
}

type User struct {
	ID       uuid.UUID
	Nickname string
}

type Message struct {
	ID        uuid.UUID
	Timestamp time.Time
	Message   string
	User      *User
}

type MessageBoard struct {
	Messages []Message
	Users    []User
}
