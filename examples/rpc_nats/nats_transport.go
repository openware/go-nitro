package main

import "github.com/statechannels/go-nitro/network/transport"

// NOTE: go-nitro do not need to implement the nats transport,
// this implementation should be provided by the user of go-nitro lib (not the lib itself).
// Exemple: go-nitro-microservice implement the nats transport, and use go-nitro lib to operate on state channels.

type natsTransport struct {
	//
}

var _ transport.Transport = (*natsTransport)(nil)

func newNatsTransport() *natsTransport {
	t := &natsTransport{
		//
	}

	return t
}

func (t *natsTransport) PollConnection() (transport.Connection, error) {
	c := newNatsConnection()

	//

	return c, nil
}

func (t *natsTransport) Close() {
	//
}
