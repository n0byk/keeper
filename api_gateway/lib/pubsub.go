package lib

import (
	"encoding/json"

	"github.com/n0byk/keeper/engine"

	log "github.com/sirupsen/logrus"
)

var (
	publish     = "publish"
	register    = "register"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
)

func (ps *engine.Pubsub) subscribe(client *engine.Client, topic string) *engine.Pubsub {
	clientSub := ps.getSubscriptions(topic, client)
	if len(clientSub) > 0 {
		// Means client is subscribed before to this topic
		return ps
	}
	newSubscription := engine.Subscription{
		Topic:  topic,
		Client: client,
	}
	ps.Subscriptions = append(ps.Subscriptions, newSubscription)
	return ps
}

func (ps *engine.Pubsub) publish(topic string, message string, excludedClient *engine.Client) {
	subscriptions := ps.getSubscriptions(topic, nil)
	for _, subscription := range subscriptions {
		log.Println("Sending to client: ", subscription.Client.ID)
		subscription.Client.Conn.WriteMessage(1, []byte(message))
		subscription.Client.send(message)
	}
}

func (client *engine.Client) send(message string) error {
	return client.Conn.WriteMessage(1, []byte(message))
}

// unsubscribe function for removing subscription from the pubsub
func (ps *engine.Pubsub) unsubscribe(client *engine.Client, topic string) *engine.Pubsub {
	for index, sub := range ps.Subscriptions {
		if sub.Client.ID == client.ID && sub.Topic == topic {
			// Found the subscription from client
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}
	return ps
}

// RemoveClient function for removing the clients from the socket
func (ps *engine.Pubsub) RemoveClient(client engine.Client) *engine.Pubsub {
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

func (ps *engine.Pubsub) getSubscriptions(topic string, client *engine.Client) []engine.Subscription {
	var subscriptionList []engine.Subscription
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
func (ps *engine.Pubsub) AddClient(client engine.Client) *engine.Pubsub {
	ps.Clients = append(ps.Clients, client)
	log.Info("Adding a new client to the list: ", client.ID)
	payload := []byte("Hello Client ID" + client.ID)
	client.Conn.WriteMessage(1, payload)
	return ps
}

// HandleReceivedMessage for hendling the recieved message from the server.
func (ps *engine.Pubsub) HandleReceivedMessage(client engine.Client, messageType int, payload []byte) *engine.Pubsub {
	m := engine.Message{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		log.Println("Something is wrong with the message!")
		return ps
	}
	switch m.Action {
	case publish:
		ps.publish(m.Topic, m.Message, nil)
		break
	case subscribe:
		ps.subscribe(&client, m.Topic)
		log.Println("New subscriber to topic: ", m.Topic, len(ps.Subscriptions))
		break
	case unsubscribe:
		ps.unsubscribe(&client, m.Topic)
		break
	case register:
		// RegistrationRequest()
		ps.unsubscribe(&client, m.Topic)
		break
	default:
		break
	}
	return ps
}
