package client

import (
	"testing"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestRpcNats(t *testing.T) {
	opts := &server.Options{}
	ns, err := server.NewServer(opts)

	assert.NoError(t, err)
	ns.Start()

	nc, err := nats.Connect(ns.ClientURL())
	println(nc)
	assert.NoError(t, err)
}
