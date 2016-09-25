package server

import (
  "log"
  "encoding/json"
  "net/http"
  "github.com/gorilla/websocket"
  "github.com/hagbarth/GoChat/app/chat"
)

var messageBoard chat.MessageBoard
var upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
}
var connections []*websocket.Conn

/**
* Responds with a JSON string from the obejct provided
* @param object {interface}
* @param w {http.ResponseWriter}
**/
func respondJSON(object interface{}, w http.ResponseWriter) {
  js, err := json.Marshal(object)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

/**
* Adds a new user and responds with the user id
**/
func handleLogin(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var u chat.User
  err := decoder.Decode(&u)
  if err != nil {
    panic("oh no!")
  }
  respondJSON(messageBoard.AddUser(u), w)
}

/**
* Handles message requests
**/
func handleMessages(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {
    addNewMessage(w, r)
    return
  }
  getMessages(w, r)
  return
}

/**
* Adds a new message and returns all messages
**/
func addNewMessage(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var m chat.MessageRequest
  err := decoder.Decode(&m)
  if err != nil {
    panic("oh no!")
  }
  messages := messageBoard.AddMessage(m.Uid, m.Message)
  respondJSON(messages, w)
}

/**
* Gets all messages
**/
func getMessages(w http.ResponseWriter, r *http.Request) {
  respondJSON(messageBoard.Messages, w)
}

/**
* Websocket server
**/
func socketHandler(w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return
  }
  connections = append(connections, conn)
  for {
    var m chat.MessageRequest
    conn.ReadJSON(&m)
    messages := messageBoard.AddMessage(m.Uid, m.Message)
    for _, connection := range connections {
      connection.WriteJSON(messages)
    }
  }
}

func Serve () {
  connections = make([]*websocket.Conn, 0)
  messageBoard = chat.NewMessageBoard()
  http.HandleFunc("/socketserver", socketHandler)
  http.Handle("/", http.FileServer(http.Dir("static")))
  http.HandleFunc("/login", handleLogin)
  http.HandleFunc("/messages", handleMessages)
  http.ListenAndServe(":8000", nil)
}
