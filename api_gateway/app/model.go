package app

import (
	"github.com/gorilla/websocket"
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
	Action    string `json:"action"`
	LobbyName string `json:"lobby_name"`
}

type TRegistrationRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TAuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TAuthenticateRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TCreateLobbyRequest struct {
	Token     string `json:"token"`
	LobbyName string `json:"lobby_name"`
}

type TAddToLobbyRequest struct {
	LobbyName   string `json:"lobby_name"`
	Token       string `json:"token"`
	PublicToken string `json:"public_token"`
}

type TGetPublicTokenRequest struct {
	Login string `json:"login"`
}

type TSubscribeToLobbyRequest struct {
	LobbyName string `json:"lobby_name"`
	Token     string `json:"token"`
}

type TSetDataRequest struct {
	Token     string `json:"token"`
	LobbyName string `json:"lobby_name"`
	LobbyId   string `json:"lobby_id"`
	RowData   string `json:"row_data"`
}

type TValidatePermissions struct {
	Token     string `json:"token"`
	LobbyName string `json:"lobby_name"`
}

type TLobbyIdRequest struct {
	LobbyName string `json:"lobby_name"`
}
