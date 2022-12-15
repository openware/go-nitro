package transport

import (
	"math/rand"
	"sync"
	"time"
)

type ChanTransport struct {
	connectionCh chan Connection
	connections  sync.Map
}

var _ Transport = (*ChanTransport)(nil)

func NewChanTransport() *ChanTransport {
	return &ChanTransport{
		connectionCh: make(chan Connection),
	}
}

func (t *ChanTransport) Connect(other *ChanTransport, minLat, maxLat time.Duration) {
	con := newChanConnection()
	otherCon := newChanConnection()

	simLat := func() {
		if minLat == 0 && maxLat == 0 {
			return
		}

		time.Sleep(minLat + time.Duration(rand.Int63n(int64(maxLat-minLat))))
	}

	go func() {
		simLat()

		t.connections.Store(con, struct{}{})

		t.connectionCh <- con

		for {
			msg, ok := <-otherCon.sendCh
			if !ok {
				break
			}

			simLat()

			con.recvCh <- msg
		}

		t.connections.Delete(con)

		close(con.recvCh)
	}()

	go func() {
		simLat()

		other.connections.Store(otherCon, struct{}{})

		other.connectionCh <- otherCon

		for {
			msg, ok := <-con.sendCh
			if !ok {
				break
			}

			simLat()

			otherCon.recvCh <- msg
		}

		other.connections.Delete(otherCon)

		close(otherCon.recvCh)
	}()
}

func (t *ChanTransport) PollConnection() (Connection, error) {
	con, ok := <-t.connectionCh
	if !ok {
		return nil, ErrTransportClosed
	}

	return con, nil
}

func (t *ChanTransport) Close() {
	close(t.connectionCh)

	t.connections.Range(func(key, value interface{}) bool {
		key.(Connection).Close()

		return true
	})
}
