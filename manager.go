package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)




var (
	webSocketUpgrader = websocket.Upgrader{
		CheckOrigin: checkOrigin,
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	Clients ClientList
	sync.Mutex

	handlers map[string] EventHandler
}

func NewManager() *Manager{
	m :=  &Manager{
		Clients: make(ClientList, 0),
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}

func(m *Manager) setupEventHandlers(){
	m.handlers[EventSendMessage] = SendMessage
	
}

func SendMessage(event Event, c *Client) error {
	fmt.Println(event)
	return nil
}

func(m *Manager) routeEvent(event Event, c *Client) error{
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil

	}else {
		return errors.New("there is no such event type")
	}
}
func(m *Manager) serveWs(w http.ResponseWriter, r *http.Request){
	log.Println("New connection")

	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("cannot upagrade connection: %v", err)
		return
	}

	client := NewClient(conn, m)

	m.AddClient(client)

	// read message
	go client.ReadMessages()
	go client.WriteMessages()
	
}

func(m *Manager) AddClient(client *Client){
	m.Lock()
	defer m.Unlock()
	
	_, ok := m.Clients[client]
	if !ok{
		m.Clients[client] = true
	}
}

func(m *Manager) RemoveClient(client *Client){
	m.Lock()
	defer m.Unlock()

	_ , ok := m.Clients[client]
	if ok {
		delete(m.Clients, client)
	}


}

func checkOrigin(r *http.Request) bool{
	origin := r.Header.Get("Origin")
	switch {
	case origin == "http://localhost:8080":
		return true
	default:
		return false
	}
}

func(m *Manager) Login(w http.ResponseWriter, r *http.Request){
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "mal formed json", http.StatusBadRequest)
		return
	}

	if input.Username != "nico" && input.Password != "123" {
		http.Error(w, "bad credenttial", http.StatusUnauthorized)
	}
}