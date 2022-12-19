package transport

import "github.com/nats-io/nats.go"

type NatsTransport struct {
	connection Connection
}

var _ Transport = (*NatsTransport)(nil)

func NewNatsTransport(connectionUrl string, pubTopicName string, subTopicName string) (*NatsTransport, error) {
	// wouldn't it be better to get messaged directly into the channel?
	subChannel := make(chan *nats.Msg)
	natsConnection, err := NewNatsConnection(connectionUrl, pubTopicName, subTopicName, subChannel)
	if err != nil {
		return nil, err
	}

	natsTransport := &NatsTransport{
		connection: natsConnection,
	}

	return natsTransport, nil
}

func (t *NatsTransport) PollConnection() (Connection, error) {
	return t.connection, nil
}

func (t *NatsTransport) Close() {
	t.connection.Close()
}
