package transport

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	netproto "github.com/statechannels/go-nitro/network/protocol"
)

type natsConnection struct {
	nc *nats.Conn

	pubTopicName string

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
		pubTopicName: pubTopicName,
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

func (c *natsConnection) Send(msg netproto.Message) {
	err := c.nc.Publish(c.pubTopicName, []byte(msg.Type()))
	if err != nil {
		log.Error().Err(err).Msgf("failed to send message on topic: %s. msg: %s", c.pubTopicName, msg)
	}
}

func (c *natsConnection) Recv() (netproto.Message, error) {
	// TODO: either store data into the natsConnection and get event by event or dirctly subscribe to channel
}

func (c *natsConnection) Close() {
	err := c.natsSubscription.Unsubscribe()
	if err != nil {
		log.Error().Err(err).Msgf("failed to unsubscribe from topic: %s.", c.subTopicName)
	}
	close(c.subChannel)
	// TODO: verify that we don't close nats connection
}
