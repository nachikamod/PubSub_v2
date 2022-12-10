package pool

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// Websocket client connection address book
type pool struct {
	Clients map[string]*websocket.Conn
}

// Message body
type Message struct {
	Topic   string          `json:"topic"`
	Message json.RawMessage `json:"message"`
}

// Construct a new pool for connections
func NewPool() *pool {
	return &pool{
		Clients: make(map[string]*websocket.Conn),
	}
}

// Add client to provided address book
func (ps *pool) AddClient(uid string, conn *websocket.Conn) {
	log.Printf("New client connected : %s, address: %s, Total clients in pool : %d", uid, conn.RemoteAddr().String(), len(ps.Clients)+1)
	ps.Clients[uid] = conn
}

// Publish message to perticular client
func (ps *pool) PublishToClient(clientId string, topic string, message json.RawMessage) {

	// Iterate over registerd clients
	for key, conn := range ps.Clients {

		// Find the client by id
		if key == clientId {

			payload, err := json.Marshal(&Message{
				Topic:   topic,
				Message: message,
			})

			if err != nil {
				log.Printf("Failed to publish message")
				return
			}

			conn.WriteMessage(1, payload)
			return
		}
	}
}

// Remove client from address book
func (ps *pool) RemoveClient(con *websocket.Conn) {
	// Iterate over registerd clients
	for key, conn := range ps.Clients {

		// Find the client by adress
		if conn.RemoteAddr().String() == con.RemoteAddr().String() {
			delete(ps.Clients, key) // remove the client
			log.Printf("Client : (%s) disconnected", key)
			break
		}
	}

}
