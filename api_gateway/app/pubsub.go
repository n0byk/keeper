package app

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var (
	publish        = "publish"
	subscribe      = "subscribe"
	unsubscribe    = "unsubscribe"
	register       = "register"
	auth           = "auth"
	LobbyCreate    = "create_lobby"
	GetPublicToken = "get_public_token"
	GetLobbyId     = "lobby_id"
	LobbyInvite    = "invite_lobby"
	setData        = "set_data"
)

func (ps *Pubsub) subscribe(client *Client, topic string) *Pubsub {
	clientSub := ps.getSubscriptions(topic, client)
	if len(clientSub) > 0 {
		// Means client is subscribed before to this topic
		return ps
	}
	newSubscription := Subscription{
		Topic:  topic,
		Client: client,
	}
	ps.Subscriptions = append(ps.Subscriptions, newSubscription)
	return ps
}

func (ps *Pubsub) publish(topic string, message string, excludedClient *Client) {
	subscriptions := ps.getSubscriptions(topic, nil)
	for _, subscription := range subscriptions {
		log.Println("Sending to client: ", subscription.Client.ID)
		subscription.Client.Conn.WriteMessage(1, []byte(message))
		subscription.Client.send(message)
	}
}

func (client *Client) send(message string) error {
	return client.Conn.WriteMessage(1, []byte(message))
}

// unsubscribe function for removing subscription from the pubsub
func (ps *Pubsub) unsubscribe(client *Client, topic string) *Pubsub {
	for index, sub := range ps.Subscriptions {
		if sub.Client.ID == client.ID && sub.Topic == topic {
			// Found the subscription from client
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}
	return ps
}

// RemoveClient function for removing the clients from the socket
func (ps *Pubsub) RemoveClient(client Client) *Pubsub {
	// First remove all the subscription of this client
	for index, sub := range ps.Subscriptions {
		if client.ID == sub.Client.ID {
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}
	// Second remove all the client from the list
	for index, c := range ps.Clients {
		if c.ID == client.ID {
			ps.Clients = append(ps.Clients[:index], ps.Clients[index+1:]...)
		}
	}
	return ps
}

func (ps *Pubsub) getSubscriptions(topic string, client *Client) []Subscription {
	var subscriptionList []Subscription
	for _, subscription := range ps.Subscriptions {
		if client != nil {
			if subscription.Client.ID == client.ID && subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		} else {
			if subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		}
	}
	return subscriptionList
}

// AddClient function for adding new client to the list
func (ps *Pubsub) AddClient(client Client) *Pubsub {
	ps.Clients = append(ps.Clients, client)
	log.Info("Adding a new client to the list: ", client.ID)
	payload := []byte("Hello anonymous")
	client.Conn.WriteMessage(1, payload)
	return ps
}

// HandleReceivedMessage for hendling the recieved message from the server.
func (ps *Pubsub) HandleReceivedMessage(client Client, messageType int, payload []byte) *Pubsub {
	m := Message{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		log.Println("Something is wrong with the message!")
		return ps
	}
	switch m.Action {
	case publish:
		message, err := SetData(payload)
		if err != nil {
			log.Println(err)
			client.send("")
			break
		}
		fmt.Println(message)
		ps.publish(m.LobbyName, message, nil)
		break
	case subscribe:
		message, err := ValidatePermissions(payload)
		if err != nil {
			log.Println(err)
			client.send("")
			break
		}
		client.send(message)
		ps.subscribe(&client, message)
		log.Println("New subscriber to topic: ", message, len(ps.Subscriptions))
		break
	case unsubscribe:
		ps.unsubscribe(&client, m.LobbyName)
		break
	case register:

		message, err := Registration(payload)
		if err != nil {
			log.Println(err)
			client.send("")
		}
		client.send(message)
		break

	case auth:
		message, err := Authenticate(payload)
		if err != nil {
			log.Println(err)
			client.send("")
		}
		client.send(message)
		break

	case LobbyCreate:
		message, err := CreateLobby(payload)
		if err != nil {
			log.Println(err)
			client.send("")
		}
		client.send(message)
		break

	case GetLobbyId:
		message, err := PublicGetLobbyId(payload)
		if err != nil {
			log.Println(err)
			client.send("")
		}
		client.send(message)
		break
	case GetPublicToken:
		message, err := PublicTokenGet(payload)
		if err != nil {
			log.Println(err)
			client.send("")
		}
		client.send(message)
		break
	case LobbyInvite:
		message, err := AddToLobby(payload)
		if err != nil {
			log.Println(err)
			client.send("")
		}
		client.send(message)
		break
	case setData:
		message, err := AddToLobby(payload)
		if err != nil {
			log.Println(err)
			client.send("")
		}
		client.send(message)
		break

	default:
		break
	}
	return ps
}
