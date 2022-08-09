package engine

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

// Pubsub structure for storing list of the clients
type Pubsub struct {
	Clients       []Client
	Subscriptions []Subscription
}

// Client structure for representing a client information
type Client struct {
	ID   string
	Conn *websocket.Conn
}

// Subscription structure denotes single subscription
type Subscription struct {
	Topic  string
	Client *Client
}

// Message structure is for storing incoming messages
type Message struct {
	Action  string `json:"action"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type Distribution struct {
	nc  *nats.Conn
	ec  *nats.EncodedConn
	ctx context.Context
}

type RegistrationRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
