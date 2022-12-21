package transport

type NatsTransport struct {
	connection Connection
}

var _ Transport = (*NatsTransport)(nil)

func NewNatsTransport(connectionUrl string, pubTopicNames []string) (*NatsTransport, error) {
	// wouldn't it be better to get messaged directly into the channel?
	natsConnection, err := NewNatsConnection(connectionUrl, pubTopicNames)
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
