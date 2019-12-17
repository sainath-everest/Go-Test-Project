package main

import (
        "fmt"
        "log"
        "net/http"
        "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan []byte)           // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{}

// Define our message object
// type Message struct {
//     Email    string `json:"e    mail"`
//     Username string `json:"username"`
//     Message  string `json:"message"`
// }
func main() {
    fmt.Println("hello world")
    // Create a simple file server
    // fs := http.FileServer(http.Dir("../public"))
    // http.Handle("/", fs)

     // Configure websocket route
     http.HandleFunc("/ws", handleConnections)

     // Start listening for incoming chat messages
     go handleMessages()

     log.Println("http server started on :8000")
     err := http.ListenAndServe(":8000", nil)
     if err != nil {
             log.Fatal("ListenAndServe: ", err)
     }
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Heyy come here")
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    // Upgrade initial GET request to a websocket
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
            log.Fatal(err)
    }
    // Make sure we close the connection when the function returns
    defer ws.Close()

    // Register our new client
    clients[ws] = true

    for {
        // var msg Message
        // Read in a new message as JSON and map it to a Message object
        _, p, err := ws.ReadMessage()
        if err != nil {
                log.Printf("error: %v", err)
                delete(clients, ws)
                break
        }
        // Send the newly received message to the broadcast channel
        broadcast <- p
}
}
func handleMessages() {
        for {
                // Grab the next message from the broadcast channel
                p := <-broadcast
                // Send it out to every client that is currently connected
                for client := range clients {
                        err := client.WriteMessage(1,p)
                        if err != nil {
                                log.Printf("error: %v", err)
                                client.Close()
                                delete(clients, client)
                        }
                }
        }
}
  