package transport

import (
	"sync"
)

type chanConnection struct {
	mu     sync.Mutex
	sendCh chan []byte
	recvCh chan []byte
}

var _ Connection = (*chanConnection)(nil)

func newChanConnection() *chanConnection {
	con := &chanConnection{
		sendCh: make(chan []byte),
		recvCh: make(chan []byte),
	}

	return con
}

func (c *chanConnection) Send(msgType string, data []byte) {
	c.sendCh <- data
}

func (c *chanConnection) Recv() ([]byte, error) {
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
