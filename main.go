package main

import "net/http"


func main(){
	
	setupApi()

	http.ListenAndServe(":8080", nil)
}

func setupApi(){
	manager := NewManager()
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWs)
	http.HandleFunc("/login", manager.Login)

}
