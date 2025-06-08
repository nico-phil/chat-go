package main

type Event struct {
	Type    string          `json:"type"`
	Payload any `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "new_message"
	EventCreateChat = "create_chat"
)

type SendMesageEvent struct {
	Message string `json:"mesage"`
	From    string `json:"from"`
}
