package main

import (
	"context"
	"net/http"
)

func main() {

	setupApi()

	http.ListenAndServe(":8080", nil)
}

func setupApi() {
	manager := NewManager()


	redisclient := NewRedisClient()

	_ = redisclient.Set(context.Background(), "hello", "world", 0)
	_ = redisclient.Set(context.Background(), "name", "nico", 0)
	

	redisclient.Sub(context.Background(), "user_A")
	
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWs)
	http.HandleFunc("/login", manager.Login)

}
