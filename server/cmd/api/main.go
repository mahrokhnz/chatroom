package main

import (
	"fmt"
	"net/http"
	"ws/pkg/ws"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Welcome to the Chat Room!")
	})
	http.HandleFunc("/ws", ws.ConnectionHandler)

	go ws.GlobalBroadcastHandler()

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
