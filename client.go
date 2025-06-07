package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	Manager    *Manager

	// egress chanel is used to avoid concurrent write to the websocket connection
	Egress chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		Manager:    manager,
		Egress:     make(chan Event),
	}
}

func (c *Client) ReadMessages() {
	//there are mutiple message type (TextMessage, BinaryMessage, CloseMessage, PingMessage and
	// PongMessage)
	defer func() {
		c.Manager.RemoveClient(c)
	}()

	err := c.connection.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		c.connection.Close()
		return
	}

	c.connection.SetReadLimit(50)

	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Println("failed to unmarshal even", err)
			break
		}

		if err := c.Manager.routeEvent(request, c); err != nil {
			log.Println("error handling mesage", err)
		}

	}
}

func (c *Client) WriteMessages() {
	defer func() {
		c.Manager.RemoveClient(c)
	}()

	ticker := time.NewTicker(pingInterval)
	for {
		select {
		case message, ok := <-c.Egress:
			if !ok {
				err := c.connection.WriteMessage(websocket.CloseMessage, nil)
				if err != nil {
					log.Printf("connection closed: %v", err)
				}
			}

			jsonData, err := json.Marshal(message)
			if err != nil {
				log.Println("error marshalling message", err)
				return
			}

			err = c.connection.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Printf("failed to send message to client: %v", err)
			}

			log.Printf("message: %s", message)

		case <-ticker.C:
			log.Println("ping")
			err := c.connection.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("failed to ping the client")
				c.connection.Close()
				return
			}
		}
	}
}

func (c *Client) pongHandler(pongMessage string) error {
	fmt.Print("pong received")

	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
