package distribution

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type Distribution struct {
	nc  *nats.Conn
	ec  *nats.EncodedConn
	ctx context.Context
}

type Registration struct {
	Login    string
	Password string
}

type Message struct {
	Action       string
	Registration Registration
}

func NewConnection() *Distribution {
	d := &Distribution{ctx: context.Background()}
	d.init()

	return d
}

// Init connects to a Nats instance
func (d *Distribution) init() {
	var err error
	if d.nc, err = nats.Connect(nats.DefaultURL); err != nil {
		panic(err)
	}
	d.ec, _ = nats.NewEncodedConn(d.nc, nats.JSON_ENCODER)
}

func (d *Distribution) Request(subj string, queue string, message *Message) Message {
	var response *Message
	err := d.ec.Request(subj, message, &response, 10*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	}
	return *response
}

func (d *Distribution) Drain() {
	d.nc.Drain()
}
