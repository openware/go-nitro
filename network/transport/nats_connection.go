package transport

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type node struct {
	prev *node
	next *node

	data *nats.Msg
}

// TODO: maybe we need mutex
type queue struct {
	head *node
	tail *node
}

func (n *queue) enqueue(data *nats.Msg) {
	newNode := &node{
		next: nil,
		prev: n.tail,
		data: data,
	}
	n.tail = newNode
}

func (n *queue) dequeue() *node {
	val := n.head
	n.head = val.next

	if val.next == nil {
		n.tail = nil
	}

	return val
}

type natsConnection struct {
	nc *nats.Conn

	subTopicName     string //TODO: convert to array of string, listen to those topics on start
	subChannel       chan *nats.Msg
	natsSubscription *nats.Subscription
	queue            *queue
}

func NewNatsConnection(connectionUrl string, subTopicName string, subChannel chan *nats.Msg) (*natsConnection, error) {
	nc, err := nats.Connect(connectionUrl)
	natsConnection := &natsConnection{
		nc:           nc,
		subTopicName: subTopicName,
		subChannel:   subChannel,
	}
	go natsConnection.handleIncomingMessages()
	err = natsConnection.subscribeWithChannel()

	return natsConnection, err
}

func (c *natsConnection) handleIncomingMessages() {
	for msg := range c.subChannel {
		c.queue.enqueue(msg)
	}
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
	msg := c.queue.dequeue()
	if msg == nil {
		if c.natsSubscription == nil {
			return nil, ErrConnectionClosed
		}

		return nil, nil
	}
	return msg.data.Data, nil
}

func (c *natsConnection) Close() {
	err := c.natsSubscription.Unsubscribe()
	if err != nil {
		log.Error().Err(err).Msgf("failed to unsubscribe from topic: %s.", c.subTopicName)
	}
	close(c.subChannel)
	// TODO: verify that we don't close nats connection
}
