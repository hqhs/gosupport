package app

import (
	"encoding/json"
	"net/http"
)

// ServeWs handles websocket requests from the peer.
func (s *Server) ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Log("err", err, "then", "During upgrading websocket connection")
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool
	// Inbound messages from the customers.
	broadcast <-chan Message
	resend chan []byte
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub(b <-chan Message) *Hub {
	return &Hub{
		broadcast:  b,
		resend:     make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case payload := <-h.resend:
			for client := range h.clients {
				select {
				case client.send <- payload:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				payload, _ := json.Marshal(message)
				select {
				case client.send <- payload:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
