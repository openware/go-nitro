package transport

import netproto "github.com/statechannels/go-nitro/network/protocol"

type Connection interface {
	Send(netproto.Message)
	Recv() (netproto.Message, error)

	Close()
}
