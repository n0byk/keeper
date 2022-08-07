package app

import (
	"fmt"
	"sync"

	"backend_service/distribution"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

func ServeMessages(conn *websocket.Conn, nc *distribution.Distribution) {
	fmt.Println("version 0.0.2")
	var mt sync.Mutex

	for {
		i, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message no.", i)
			conn.Close()
			return
		}

		// un/subscribe if event == un/subscribe.
		var smsg = string(msg)
		action := gjson.Get(smsg, "action").String()
		channel := gjson.Get(smsg, "channel").String()
		data := gjson.Get(smsg, "data").String()
		password := gjson.Get(smsg, "password").String()
		login := gjson.Get(smsg, "login").String()

		switch action {

		case "registration":
			message := &distribution.Message{Action: action, Registration: distribution.Registration{login, password}}

			response := nc.Request("keeper", "keeper", message)
			Publish(i, channel, []byte(response))

		case "message":
			message := &distribution.Message{Action: action}

			nc.Request("keeper", "keeper", message)
			Publish(i, channel, []byte(data))

		case "subscribe":

			Subscribe(channel, conn)
			msg = []byte("subscribe to " + channel + " success!")

		case "unsubscribe":

			Unsubscribe(channel, conn)
			msg = []byte("unsubscribe from " + channel + " success!")
		default:
			continue

		}

		fmt.Println(string(msg))

		mt.Lock()
		if err = conn.WriteMessage(i, msg); err != nil {
			fmt.Println(err)
		}
		mt.Unlock()
	}
}
