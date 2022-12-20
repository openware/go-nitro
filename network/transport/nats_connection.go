package transport

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type natsConnection struct {
	nc *nats.Conn

	subTopicName     string
	subChannel       chan *nats.Msg
	natsSubscription *nats.Subscription
}

var _ Connection = (*natsConnection)(nil)

func NewNatsConnection(connectionUrl string, pubTopicName string, subTopicName string, subChannel chan *nats.Msg) (*natsConnection, error) {
	nc, err := nats.Connect(connectionUrl)
	natsConnection := &natsConnection{
		nc:           nc,
		subTopicName: subTopicName,
		subChannel:   subChannel,
	}

	err = natsConnection.subscribeWithChannel()

	return natsConnection, err
}

func (c *natsConnection) subscribeWithChannel() error {
	sub, err := c.nc.ChanSubscribe(c.subTopicName, c.subChannel)
	if err != nil {
		return err
	}

	c.natsSubscription = sub
	return nil
}

func (c *natsConnection) Send(t string, data []byte) {
	err := c.nc.Publish(t, data)
	if err != nil {
		log.Error().Err(err).Msgf("failed to send message on topic: %s. msg: %s", t, string(data))
	}
}

func (c *natsConnection) Recv() ([]byte, error) {
	// TODO: either store data into the natsConnection and get event by event or dirctly subscribe to channel
	return nil, nil
}

func (c *natsConnection) Close() {
	err := c.natsSubscription.Unsubscribe()
	if err != nil {
		log.Error().Err(err).Msgf("failed to unsubscribe from topic: %s.", c.subTopicName)
	}
	close(c.subChannel)
	// TODO: verify that we don't close nats connection
}
