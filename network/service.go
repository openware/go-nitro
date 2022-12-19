package network

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/internal"
	"github.com/statechannels/go-nitro/network/transport"
)

// TODO: replace Poll methods (sub packages too) with data channel + error channel

type Service struct {
	Logger zerolog.Logger

	Transport transport.Transport
}

func NewService(trp transport.Transport) *Service {
	s := &Service{
		Transport: trp,
	}

	return s
}

func (s *Service) PollPeer() (*NetworkServiceConnection, error) {
	con, err := s.Transport.PollConnection()
	if err != nil {
		if errors.Is(err, transport.ErrTransportClosed) {
			s.Logger.Info().Msg("transport closed")

			return nil, internal.WrapError(ErrServiceClosed, err)
		}

		return nil, err
	}

	p := NewNetworkServiceConnection(uuid.New(), con)
	p.Logger = s.Logger.With().Str("peer", p.Id.String()).Logger()

	s.Logger.Info().Str("peer", p.Id.String()).Msg("peer connected")

	return p, nil
}

func (s *Service) Close() {
	s.Transport.Close()
}
