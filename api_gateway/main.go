package main

import (
	"net/http"

	"github.com/n0byk/keeper/engine"

	"github.com/n0byk/keeper/api_gateway/lib/distribution"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// IDGen for generating unique ids
func IDGen() string {
	return uuid.Must(uuid.New(), nil).String()
}

var ps = &engine.Pubsub{}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := engine.Client{
		ID:   IDGen(),
		Conn: conn,
	}
	// Have to add this client into the list of the clients
	ps.AddClient(client)
	log.Println("A new client is connected, Total: ", len(ps.Clients))

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			// This client is disconnected or there is error in connection
			// Have to remove/unsubscribe this client from pubsub
			ps.RemoveClient(client)
			log.Println("Total clients and subscriptions ", len(ps.Clients), len(ps.Subscriptions))
			return
		}
		ps.HandleReceivedMessage(client, messageType, p)
	}
}

func main() {
	distribution.NewConnection()
	http.HandleFunc("/ws", websocketHandler)
	log.Info("Server started on port 8080...")
	go http.ListenAndServe(":8080", nil)
}
