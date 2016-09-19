package chat

import (
  "github.com/satori/go.uuid"
  "time"
)

func NewMessageBoard() MessageBoard {
  return MessageBoard{Users: make([]User, 0), Messages: make([]Message, 0)}
}

func (m *MessageBoard) AddUser(u User) User {
  u.ID = uuid.NewV4()
  m.Users = append(m.Users, u)
  return u
}

func (m *MessageBoard) AddMessage(uid string, mes string) []Message {
  userUuid, err := uuid.FromString(uid)
  if err != nil {
    panic("User id conversion error")
  }
  message := Message{
    ID:        uuid.NewV4(),
    Timestamp: time.Now(),
    Message:   mes,
  }
  for _, u := range m.Users {
    if uuid.Equal(u.ID, userUuid) {
      message.User = &u
    }
  }
  m.Messages = append(m.Messages, message)
  return m.Messages
}
