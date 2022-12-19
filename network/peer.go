package network

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/internal"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/network/transport"
)

// TODO: use mutexes for both maps and logs (removing sync.Map usage)

type NetworkServiceConnection struct {
	Logger zerolog.Logger

	Id uuid.UUID

	Connection transport.Connection

	errorHandlers   sync.Map
	messageHandlers sync.Map

	requests       sync.Map
	responseErrChs sync.Map
	responseChs    sync.Map
}

func NewNetworkServiceConnection(id uuid.UUID, con transport.Connection) *NetworkServiceConnection {
	p := &NetworkServiceConnection{
		Id: id,

		Connection: con,
	}

	go p.handleMessages()

	return p
}

func RegisterErrorHandler[T netproto.Message](p *NetworkServiceConnection, handle func(*netproto.Error, T)) {
	var msg T

	p.errorHandlers.Store(msg.Type(), func(errMsg *netproto.ErrorMessage) {
		msg, ok := errMsg.Message.(T)
		if !ok {
			p.SendError(netproto.ErrInvalidMessage, errMsg)

			return
		}

		handle(errMsg.Error, msg)
	})

	p.Logger.Trace().Str("message_type", msg.Type()).Msg("registered error handler")
}

func RegisterErrorHandlerOnce[T netproto.Message](p *NetworkServiceConnection, handle func(*netproto.Error, T)) {
	RegisterErrorHandler(p, func(err *netproto.Error, msg T) {
		UnregisterErrorHandler[T](p)

		handle(err, msg)
	})
}

func UnregisterErrorHandler[T netproto.Message](p *NetworkServiceConnection) {
	var msg T

	p.errorHandlers.Delete(msg.Type())

	p.Logger.Trace().Str("message_type", msg.Type()).Msg("unregistered error handler")
}

func RegisterMessageHandler[T netproto.Message](
	p *NetworkServiceConnection,
	handle func(T) *netproto.Error,
	handleErr func(*netproto.Error, T),
) {
	var msg T

	p.messageHandlers.Store(msg.Type(), func(m netproto.Message) {
		msg, ok := m.(T)
		if !ok {
			p.SendError(netproto.ErrInvalidMessage, msg)

			return
		}

		err := handle(msg)
		if err != nil {
			if handleErr != nil {
				handleErr(err, msg)
			}

			p.SendError(err, msg)

			return
		}
	})

	p.Logger.Trace().Str("message_type", msg.Type()).Msg("registered message handler")
}

func RegisterMessageHandlerOnce[T netproto.Message](
	p *NetworkServiceConnection,
	handle func(T) *netproto.Error,
	handleErr func(*netproto.Error, T),
) {
	RegisterMessageHandler(
		p,
		func(msg T) *netproto.Error {
			err := handle(msg)

			UnregisterMessageHandler[T](p)

			return err
		},
		handleErr,
	)
}

func UnregisterMessageHandler[T netproto.Message](p *NetworkServiceConnection) {
	var msg T

	p.messageHandlers.Delete(msg.Type())

	p.Logger.Trace().Str("message_type", msg.Type()).Msg("unregistered message handler")
}

func RegisterRequestHandler[
	Req netproto.RequestMessage,
	Res netproto.ResponseMessage,
](p *NetworkServiceConnection, handle func(Req) (Res, *netproto.Error), handleErr func(*netproto.Error, Req)) {
	RegisterMessageHandler(
		p,
		func(req Req) *netproto.Error {
			res, err := handle(req)
			if err != nil {
				return err
			}

			p.SendResponse(req.Id(), res)

			return nil
		},
		handleErr,
	)
}

func RegisterRequestHandlerOnce[
	Req netproto.RequestMessage,
	Res netproto.ResponseMessage,
](p *NetworkServiceConnection, handle func(Req) (Res, *netproto.Error), handleErr func(*netproto.Error, Req)) {
	RegisterRequestHandler(
		p,
		func(req Req) (Res, *netproto.Error) {
			res, err := handle(req)

			UnregisterMessageHandler[Req](p)

			return res, err
		},
		handleErr,
	)
}

func UnregisterRequestHandler[Req netproto.RequestMessage](p *NetworkServiceConnection) {
	UnregisterMessageHandler[Req](p)
}

func RegisterResponseHandler[Res netproto.ResponseMessage, Req netproto.RequestMessage](
	p *NetworkServiceConnection,
	handleReqErr func(*netproto.Error, Req),
	handle func(Res, Req) *netproto.Error,
) {
	if handleReqErr != nil {
		RegisterErrorHandler(p, handleReqErr)
	}

	RegisterMessageHandler(
		p,
		func(res Res) *netproto.Error {
			req, ok := GetRequest[Req](p, res.RequestId())
			if !ok {
				return netproto.ErrInvalidResponse
			}

			return handle(res, req)
		},
		func(err *netproto.Error, res Res) {
			p.failRequest(res.RequestId(), internal.WrapError(ErrResponseError, err))
		},
	)
}

func RegisterResponseHandlerOnce[
	Res netproto.ResponseMessage,
	Req netproto.RequestMessage,
](p *NetworkServiceConnection, handleReqErr func(*netproto.Error, Req), handle func(Res, Req) *netproto.Error) {
	RegisterResponseHandler(
		p,
		func(err *netproto.Error, req Req) {
			UnregisterResponseHandler[Res, Req](p)

			if handleReqErr != nil {
				handleReqErr(err, req)
			}
		},
		func(res Res, req Req) *netproto.Error {
			UnregisterResponseHandler[Res, Req](p)

			return handle(res, req)
		},
	)
}

func UnregisterResponseHandler[
	Res netproto.ResponseMessage,
	Req netproto.RequestMessage,
](p *NetworkServiceConnection) {
	UnregisterMessageHandler[Res](p)
	UnregisterErrorHandler[Req](p)
}

func (p *NetworkServiceConnection) handleMessages() {
	for {
		msg, err := p.Connection.Recv()
		if err != nil {
			if errors.Is(err, transport.ErrConnectionClosed) {
				p.Logger.Info().Msg("connection closed")

				p.dropRequests()

				break
			}

			// TODO: handle error
			p.Logger.Fatal().Err(err).Msg("failed to receive message")
		}

		// NOTE: we do not hande messages in a separate goroutine
		// to ensure that messages are handled in the order they are received
		// and to avoid inconsistencies in the state of the peer
		p.handleMessage(msg, 0)
	}
}

func (p *NetworkServiceConnection) SendMessage(msg netproto.Message) {
	p.Connection.Send(msg)

	p.Logger.Trace().
		Str("message_type", msg.Type()).
		Interface("message_data", msg).
		Msg("sent message")
}

func (p *NetworkServiceConnection) SendError(err *netproto.Error, msg netproto.Message) {
	p.SendMessage(&netproto.ErrorMessage{
		Error:   err,
		Message: msg,
	})
}

func (p *NetworkServiceConnection) SendRequest(req netproto.RequestMessage) (netproto.ResponseMessage, error) {
	req.SetId(uuid.New())

	rid := req.Id()

	p.requests.Store(rid, req)
	defer p.requests.Delete(rid)

	errCh := make(chan error)
	p.responseErrChs.Store(rid, errCh)
	defer p.responseErrChs.Delete(rid)

	resCh := make(chan netproto.ResponseMessage)
	p.responseChs.Store(rid, resCh)
	defer p.responseChs.Delete(rid)

	p.SendMessage(req)

	select {
	case err := <-errCh:
		return nil, err

	case res := <-resCh:
		return res, nil
	}
}

func SendRequest[
	Req netproto.RequestMessage,
	Res netproto.ResponseMessage,
](p *NetworkServiceConnection, req Req) (Res, error) {
	res, err := p.SendRequest(req)
	if err != nil {
		return *new(Res), err
	}

	resp, ok := res.(Res)
	if !ok {
		if _, ok := p.messageHandlers.Load(res.Type()); !ok {
			p.SendError(netproto.ErrInvalidResponse, res)
		}

		return *new(Res), internal.WrapError(ErrResponseError, netproto.ErrInvalidResponse)
	}

	return resp, nil
}

func SendRequestWithHandler[
	Req netproto.RequestMessage,
	Res netproto.ResponseMessage,
](
	p *NetworkServiceConnection,
	req Req,
	handleErr func(*netproto.Error, Req),
	handle func(Res, Req) *netproto.Error,
) (Res, error) {
	RegisterResponseHandlerOnce(p, handleErr, handle)

	return SendRequest[Req, Res](p, req)
}

func (p *NetworkServiceConnection) GetRequest(rid netproto.RequestId) (netproto.RequestMessage, bool) {
	req, ok := p.requests.Load(rid)
	if !ok {
		return nil, false
	}

	return req.(netproto.RequestMessage), true
}

func GetRequest[Req netproto.RequestMessage](p *NetworkServiceConnection, rid netproto.RequestId) (Req, bool) {
	r, ok := p.GetRequest(rid)
	if !ok {
		return *new(Req), false
	}

	req, ok := r.(Req)
	if !ok {
		return *new(Req), false
	}

	return req, true
}

func (p *NetworkServiceConnection) failRequest(rid netproto.RequestId, err error) bool {
	errCh, ok := p.responseErrChs.Load(rid)
	if ok {
		p.responseChs.Delete(rid)

		errCh.(chan error) <- err
	}

	return ok
}

func (p *NetworkServiceConnection) completeRequest(rid netproto.RequestId, res netproto.ResponseMessage) bool {
	resCh, ok := p.responseChs.Load(rid)
	if ok {
		resCh.(chan netproto.ResponseMessage) <- res
	}

	return ok
}

func (p *NetworkServiceConnection) dropRequests() {
	p.responseErrChs.Range(func(key, value interface{}) bool {
		value.(chan error) <- ErrPeerClosed

		return true
	})
}

func (p *NetworkServiceConnection) SendResponse(rid netproto.RequestId, res netproto.ResponseMessage) {
	res.SetRequestId(rid)

	p.SendMessage(res)
}

func (p *NetworkServiceConnection) handleMessage(msg netproto.Message, retryCount uint) {
	p.Logger.Trace().
		Str("message_type", msg.Type()).
		Interface("message_data", msg).
		Msg("received message")

	handler, ok := p.messageHandlers.Load(msg.Type())

	switch msg := msg.(type) {
	case *netproto.ErrorMessage:
		if req, ok := msg.Message.(netproto.RequestMessage); ok {
			defer p.failRequest(req.Id(), internal.WrapError(ErrRequestError, msg.Error))
		}

		errHandler, ok := p.errorHandlers.Load(msg.Message.Type())
		if ok {
			errHandler.(func(*netproto.ErrorMessage))(msg)

			return
		}

		if req, ok := msg.Message.(netproto.RequestMessage); ok {
			if _, ok := p.responseErrChs.Load(req.Id()); ok {
				return
			}
		}

	case netproto.ResponseMessage:
		if _, ok := p.responseChs.Load(msg.RequestId()); !ok {
			p.SendError(netproto.ErrUnexpectedResponse, msg)

			return
		}

		defer p.completeRequest(msg.RequestId(), msg)

	default:
		if !ok {
			if retryCount < 3 {
				time.Sleep(time.Millisecond * 100)

				p.handleMessage(msg, retryCount+1)

				return
			}

			p.SendError(netproto.ErrUnexpectedMessage, msg)

			return
		}
	}

	if ok {
		handler.(func(netproto.Message))(msg)
	}
}

func (p *NetworkServiceConnection) Close() {
	p.Connection.Close()
}
