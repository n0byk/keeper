package distribution

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type Distribution struct {
	nc  *nats.Conn
	ec  *nats.EncodedConn
	ctx context.Context
}

func NewConnection() *Distribution {
	d := &Distribution{ctx: context.Background()}
	d.init()

	return d
}

type Registration struct {
	Login    string
	Password string
}

type Message struct {
	Action       string
	Registration Registration
}

// Init connects to a Nats instance
func (d *Distribution) init() {
	var err error
	if d.nc, err = nats.Connect(nats.DefaultURL); err != nil {
		panic(err)
	}
	d.ec, _ = nats.NewEncodedConn(d.nc, nats.JSON_ENCODER)
}

func (d *Distribution) InitSubscribe(subj string, queue string) {

	d.ec.Subscribe(subj, func(subj, reply string, msg *Message) {
		fmt.Printf("Received a person: %+v\n", msg)

		data, err := ActionExchange(*msg)
		if err != nil {
			log.Println(err)
			d.ec.Publish(reply, err)

		}
		log.Println(data)
		d.ec.Publish(reply, data)
	})

}

func (d *Distribution) Drain() {
	d.nc.Drain()
}
