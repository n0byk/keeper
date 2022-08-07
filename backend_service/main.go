package main

import (
	"backend_service/app"
	"backend_service/distribution"
	"net/http"

	"github.com/gorilla/websocket"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 2, 2) //1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", 404)
	}
	nc := distribution.NewConnection()

	go app.ServeMessages(conn, nc)
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	panic(http.ListenAndServe(":8080", nil))
}
