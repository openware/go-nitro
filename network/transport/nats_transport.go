package transport

import "github.com/nats-io/nats.go"

type NatsTransport struct {
	connection Connection
}

var _ Transport = (*NatsTransport)(nil)

func NewNatsTransport(nc *nats.Conn, pubTopicNames []string) (*NatsTransport, error) {
	// wouldn't it be better to get messaged directly into the channel?
	connection := NewNatsConnection(nc, pubTopicNames)
	natsTransport := &NatsTransport{
		connection: connection,
	}

	return natsTransport, nil
}

func (t *NatsTransport) PollConnection() (Connection, error) {
	return t.connection, nil
}

func (t *NatsTransport) Close() {
	t.connection.Close()
}
