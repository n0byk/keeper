package app

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	MainQueue = "keeper.*"

	registration = "keeper.registration"
)

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

func (d *Distribution) RegistrationRequest(data *RegistrationRequest) string {
	var response string
	err := d.ec.Request(registration, data, response, 10*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	}
	return response
}

func (d *Distribution) Drain() {
	d.nc.Drain()
}
