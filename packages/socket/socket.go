package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Global variable to store active WebSocket connections
var Clients = make(map[*websocket.Conn]bool)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func StartWebSocketServer() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/ws", wsHandler)
	log.Println("WebSocket server starting on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade error:", err)
		return
	}
	// defer conn.Close()

	// Register new client
	Clients[conn] = true

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			delete(Clients, conn)
			break
		}

		log.Printf("Received raw message: %s\n", string(message))
		handleGameMessage(message)
	}
}

func handleGameMessage(message []byte) {
	fmt.Println(message)
	// Define a struct to unmarshal your JSON data (adjust fields as needed)
	type gameCommand struct {
		Command string `json:"command"`
		Data    string `json:"data"`
	}

	var cmd gameCommand

	// Unmarshal the JSON data into the struct

	if err := json.Unmarshal([]byte(message), &cmd); err != nil {
		log.Printf("Error unmarshalling message: %s", err)
		return
	}

	// Handle the command
	fmt.Println(cmd)
	switch cmd.Command {
	case "move":
		// Handle move command
		log.Printf("Move command received with data: %s", cmd.Data)
		// Implement the logic to handle move command here
	case "action":
		// Handle action command
		log.Printf("Action command received with data: %s", cmd.Data)
		// Implement the logic to handle action command here
	default:
		log.Printf("Unknown command: %s", cmd.Command)
	}

	// You can also broadcast messages to all connected clients
	for client := range Clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte("some response")); err != nil {
			log.Printf("Error writing message: %s", err)
			client.Close()
			delete(Clients, client)
		}
	}
}
