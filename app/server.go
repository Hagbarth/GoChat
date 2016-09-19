package main

import (
	"encoding/json"
	"github.com/hagbarth/GoChat/app/chat"
	"net/http"
)

var messageBoard chat.MessageBoard

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

func main() {
	messageBoard = chat.NewMessageBoard()
	http.Handle("/", http.FileServer(http.Dir("../static")))
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/messages", handleMessages)
	http.ListenAndServe(":8000", nil)
}
