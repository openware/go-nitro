package network

import (
	"errors"
	"sync"

	"github.com/rs/zerolog"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/network/serde"
	"github.com/statechannels/go-nitro/network/transport"
)

type NetworkService struct {
	Logger     zerolog.Logger
	Connection transport.Connection
	Serde      serde.Serde

	handlerRequest      sync.Map
	handlerResponse     sync.Map
	handlerError        sync.Map
	handlerPublicEvent  sync.Map
	handlerPrivateEvent sync.Map
}

func NewNetworkService(con transport.Connection, srd serde.Serde) *NetworkService {
	p := &NetworkService{
		Connection: con,
		Serde:      srd,
	}

	go p.handleMessages()

	return p
}

func (p *NetworkService) RegisterRequestHandler(method string, handler func(*netproto.Message)) {
	p.handlerRequest.Store(method, handler)
	p.Logger.Trace().Str("method", method).Msg("registered request handler")
}

func (p *NetworkService) UnregisterRequestHandler(method string) {
	p.handlerRequest.Delete(method)
	p.Logger.Trace().Str("method", method).Msg("unregistered request handler")
}

func (p *NetworkService) RegisterErrorHandler(method string, handler func(*netproto.Message)) {
	p.handlerError.Store(method, handler)
	p.Logger.Trace().Str("method", method).Msg("registered error handler")
}

func (p *NetworkService) UnregisterErrorHandler(method string) {
	p.handlerError.Delete(method)
	p.Logger.Trace().Str("method", method).Msg("unregistered error handler")
}

func (p *NetworkService) RegisterResponseHandler(method string, handler func(*netproto.Message)) {
	p.handlerResponse.Store(method, handler)
	p.Logger.Trace().Str("method", method).Msg("registered response handler")
}

func (p *NetworkService) UnregisterResponseHandler(method string) {
	p.handlerResponse.Delete(method)
	p.Logger.Trace().Str("method", method).Msg("unregistered response handler")
}

// TODO: implement (un)registerPublicEventHandler
// TODO: implement (un)registerPrivateEventHandler

func (p *NetworkService) handleMessages() {
	for {
		data, err := p.Connection.Recv()
		if err != nil {
			if errors.Is(err, transport.ErrConnectionClosed) {
				p.Logger.Info().Msg("connection closed")
				break
			}

			// TODO: handle error
			p.Logger.Fatal().Err(err).Msg("failed to receive message")
		}

		msg, err := p.Serde.Deserialize(data)

		if err != nil {
			p.Logger.Error().Err(err).Msg("failed to deserialize message")
			return
		}

		// NOTE: we do not hande messages in a separate goroutine
		// to ensure that messages are handled in the order they are received
		// and to avoid inconsistencies in the state of the peer
		p.handleMessage(msg)
	}
}

func (p *NetworkService) SendMessage(msg *netproto.Message) {
	data, err := p.Serde.Serialize(msg)

	if err != nil {
		// TODO: handle error
		p.Logger.Error().Err(err).Msg("failed to serialize message")
		return
	}

	p.Connection.Send(msg.Method, data)

	p.Logger.Trace().
		Str("msg_type", netproto.TypeStr(msg.Type)).
		Str("method", msg.Method).
		Msg("sent message")
}

// Maybe return func(*netproto.Message). if it doesn't exist it will just return nil
func (p *NetworkService) getHandler(msg *netproto.Message) (any, bool) {
	switch msg.Type {
	case netproto.TypeRequest:
		return p.handlerRequest.Load(msg.Method)

	case netproto.TypeResponse:
		return p.handlerResponse.Load(msg.Method)

	case netproto.TypeError:
		return p.handlerError.Load(msg.Method)

	}
	// TODO: case handlerPublicEvent
	// TODO: case handlerPrivateEvent

	return nil, false
}

// TODO: I feel like this should return error
func (p *NetworkService) handleMessage(msg *netproto.Message) {
	p.Logger.Trace().
		Str("msg_type", netproto.TypeStr(msg.Type)).
		Str("method", msg.Method).
		Msg("received message")

	h, ok := p.getHandler(msg)

	if !ok {
		p.Logger.Error().
			Str("msg_type", netproto.TypeStr(msg.Type)).
			Str("method", msg.Method).
			Msg("missing handler")
		return
	}

	h.(func(*netproto.Message))(msg)
}

func (p *NetworkService) Close() {
	p.Connection.Close()
}
