package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type message struct {
	ID   string
	Data string
}

type client struct {
	conn *websocket.Conn
	mu   sync.Mutex
	ID   string
	hub  *hub
}

type hub struct {
	clients map[string]*client

	// Inbound messages from the clients.
	send chan message

	// Register requests from the clients.
	register chan *client

	// Unregister requests from clients.
	unregister chan *client

	connections map[string]*client
}

func newHub() *hub {
	return &hub{
		send:       make(chan message),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[string]*client),
	}
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)

			}
		case message := <-h.send:
			if client, ok := h.clients[message.ID]; ok {
				fmt.Println(message, "hihow")
				err := client.conn.WriteJSON(message)
				if err != nil {
					log.Printf("error: %v", err)
					client.conn.Close()
				}
			}
		}
	}
}

// Configure the upgrader
var upgrader = websocket.Upgrader{}

func main() {

	// Configure websocket route
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(hub, w, r)
	})

	// Start listening for incoming chat messages

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func handleConnections(hub *hub, w http.ResponseWriter, r *http.Request) {
	fmt.Println("i am new connection ")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ids, ok := r.URL.Query()["id"]

	if !ok {
		log.Println("Url Param 'key' is missing")
		return
	}
	id := ids[0]
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	// defer ws.Close()

	client := &client{ID: id, hub: hub, conn: ws}

	// Register our new client
	client.hub.register <- client

	for {
		var msg message
		// Read in a new message as JSON and map it to a Message object
		err := client.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			// delete(hub.connections, client)
			break
		}
		// Send the newly received message to the broadcast channel
		client.hub.send <- msg
	}

}
