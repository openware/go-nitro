package rpc

import (
	"errors"

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

	Client         *client.Client
	NetworkService *network.Service
}

func NewService(cli *client.Client, nts *network.Service) *Service {
	s := &Service{
		Client:         cli,
		NetworkService: nts,
	}

	go func() {
		for {
			peer, err := s.NetworkService.PollPeer()
			if err != nil {
				if errors.Is(err, network.ErrServiceClosed) {
					s.Logger.Info().Msg("network service closed")

					break
				}

				// TODO: handle error
				s.Logger.Fatal().Err(err).Msg("failed to poll peer")
			}

			go s.HandlePeer(peer)
		}
	}()

	return s
}

func (s *Service) HandlePeer(p *network.Peer) {
	network.RegisterRequestHandler(
		p,
		func(req *rpcproto.DirectFundRequest) (*rpcproto.DirectFundResponse, *netproto.Error) {
			s.Client.Engine.ObjectiveRequestsFromAPI <- req.Params

			res := &rpcproto.DirectFundResponse{
				Result: req.Params.Response(*s.Client.Address, s.Client.ChainId),
			}

			return res, nil
		},
		nil,
	)

	go func() {
		// TODO: handle peer disconnect

		for {
			select {
			case oid := <-s.Client.FailedObjectives():
				p.SendMessage(
					&rpcproto.ObjectiveFailed{
						ObjectiveId: oid,
						Reason:      "unknown reason", // TODO: get reason
					},
				)

			// TODO: hook into engine to send intermediate objective statuses
			case oid := <-s.Client.CompletedObjectives():
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

func (s *Service) Close() {
	s.NetworkService.Close()
}
