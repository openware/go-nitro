package main

import (
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/network/transport"
)

type natsConnection struct {
	//
}

var _ transport.Connection = (*natsConnection)(nil)

func newNatsConnection() *natsConnection {
	c := &natsConnection{
		//
	}

	return c
}

func (c *natsConnection) Send(msg netproto.Message) {
	// TODO: encode

	//
}

func (c *natsConnection) Recv() (netproto.Message, error) {
	// TODO: decode

	//

	return nil, nil
}

func (c *natsConnection) Close() {
	//
}
