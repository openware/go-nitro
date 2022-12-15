package transport

type NatsTransport struct {
	//
}

var _ Transport = (*NatsTransport)(nil)

func NewNatsTransport() *NatsTransport {
	t := &NatsTransport{
		//
	}

	return t
}

func (t *NatsTransport) PollConnection() (Connection, error) {
	c := newNatsConnection()

	//

	return c, nil
}

func (t *NatsTransport) Close() {
	//
}
