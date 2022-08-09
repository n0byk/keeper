package distribution

import (
	"context"
	"fmt"
	"time"

	"github.com/n0byk/keeper/engine"

	"github.com/nats-io/nats.go"
)

func NewConnection() *engine.Distribution {
	d := &engine.Distribution{ctx: context.Background()}
	d.init()

	return d
}

// Init connects to a Nats instance
func (d *engine.Distribution) init() {
	var err error
	if d.nc, err = nats.Connect(nats.DefaultURL); err != nil {
		panic(err)
	}
	d.ec, _ = nats.NewEncodedConn(d.nc, nats.JSON_ENCODER)
}

func (d *engine.Distribution) RegistrationRequest(data *engine.RegistrationRequest) string {
	var response string
	err := d.ec.Request(engine.Registration, data, response, 10*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	}
	return response
}

func (d *engine.Distribution) Drain() {
	d.nc.Drain()
}
