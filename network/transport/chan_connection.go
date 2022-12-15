package transport

import (
	"sync"

	netproto "github.com/statechannels/go-nitro/network/protocol"
)

type chanConnection struct {
	mu     sync.Mutex
	sendCh chan netproto.Message
	recvCh chan netproto.Message
}

var _ Connection = (*chanConnection)(nil)

func newChanConnection() *chanConnection {
	con := &chanConnection{
		sendCh: make(chan netproto.Message),
		recvCh: make(chan netproto.Message),
	}

	return con
}

func (c *chanConnection) Send(msg netproto.Message) {
	c.sendCh <- msg
}

func (c *chanConnection) Recv() (netproto.Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	msg, ok := <-c.recvCh
	if !ok {
		return msg, ErrConnectionClosed
	}

	return msg, nil
}

func (c *chanConnection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	close(c.sendCh)
	close(c.recvCh)
}
