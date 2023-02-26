package main

import (
	"github.com/gorilla/websocket"
	"log"
)

// Message Define a message struct for incoming and outgoing messages
type Message struct {
	Text string `json:"text"`
}

// Response Define a response struct for incoming and outgoing responses
type Response struct {
	Text string `json:"text"`
}

// Define a WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (app *Config) chatWebsocket(conn *websocket.Conn) {
	defer conn.Close()

	for {
		// Read a message from the WebSocket
		var message Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Create a response message
		response := Response{Text: "中医小聪" + message.Text}
		// Send the response message over the WebSocket
		err = conn.WriteJSON(response)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
