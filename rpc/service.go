package rpc

import (
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/network"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/protocols"
	rpcproto "github.com/statechannels/go-nitro/rpc/protocol"
)

// TODO: logging

type Service struct {
	Logger zerolog.Logger

	cli *client.Client
}

func NewService(cli *client.Client) *Service {
	s := &Service{
		cli: cli,
	}

	return s
}

func (s *Service) HandlePeer(p *network.Peer) {
	network.RegisterRequestHandler(
		p,
		func(req *rpcproto.DirectFundRequest) (*rpcproto.DirectFundResponse, *netproto.Error) {
			s.cli.Engine.ObjectiveRequestsFromAPI <- req.Params

			res := &rpcproto.DirectFundResponse{
				Result: req.Params.Response(*s.cli.Address, s.cli.ChainId),
			}

			return res, nil
		},
		nil,
	)

	go func() {
		// TODO: handle peer disconnect

		for {
			select {
			case oid := <-s.cli.FailedObjectives():
				p.SendMessage(
					&rpcproto.ObjectiveFailed{
						ObjectiveId: oid,
						Reason:      "unknown reason", // TODO: get reason
					},
				)

			// TODO: hook into engine to send intermediate objective statuses
			case oid := <-s.cli.CompletedObjectives():
				p.SendMessage(
					&rpcproto.ObjectiveUpdated{
						ObjectiveId: oid,
						Status:      protocols.Completed,
					},
				)
			}
		}
	}()
}
