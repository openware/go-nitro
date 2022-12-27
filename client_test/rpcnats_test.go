package client

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/network"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/network/serde"
	"github.com/statechannels/go-nitro/network/transport"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	rpcproto "github.com/statechannels/go-nitro/rpc/protocol"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/assert"
)

func setupClientWithP2PMessageService(pk []byte, port int, chain *chainservice.MockChainService, logDestination io.Writer) (client.Client, *p2pms.P2PMessageService) {
	messageservice := p2pms.NewMessageService("127.0.0.1", port, pk)
	storeA := store.NewMemStore(pk)

	return client.New(messageservice, chain, storeA, logDestination, &engine.PermissivePolicy{}, nil), messageservice
}

// User real blockchain simulated_backend_service
func TestRunRpcNats(t *testing.T) {
	wg := sync.WaitGroup{}
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:           os.Stdout,
		TimeFormat:    time.RFC3339,
		PartsOrder:    []string{"time", "level", "caller", "client", "scope", "message"},
		FieldsExclude: []string{"time", "level", "caller", "message", "client", "scope"},
	}).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("client", "").
		Str("scope", "").
		Logger()

	opts := &server.Options{}
	ns, err := server.NewServer(opts)

	assert.NoError(t, err)
	ns.Start()

	nc, err := nats.Connect(ns.ClientURL())
	chain := chainservice.NewMockChain()

	alice := testactors.Alice
	bob := testactors.Bob

	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	// TODO: if we setup chain service how do we connect everything to engine
	//sim, bindings, ethAccounts, err := chainservice.SetupSimulatedBackend(3)
	//if err != nil {
	//	t.Fatal(err)
	//}

	trp := transport.NewNatsTransport(nc, []string{fmt.Sprintf("nitro.%s", rpcproto.DirectFundRequestMethod), "nitro.test-topic"})

	clientA, msgA := setupClientWithP2PMessageService(alice.PrivateKey, 3005, chainServiceA, logger)
	clientB, msgB := setupClientWithP2PMessageService(bob.PrivateKey, 3006, chainServiceB, logger)

	con, err := trp.PollConnection()
	if err != nil {
		assert.NoError(t, err)
	}

	nts := network.NewNetworkService(con, &serde.MsgPack{})
	nts.RegisterResponseHandler(rpcproto.DirectFundRequestMethod, func(m *netproto.Message) {
		logger.Info().Msgf("Objective updated: %v", *m)
	})

	nts.RegisterErrorHandler(rpcproto.DirectFundRequestMethod, func(m *netproto.Message) {
		logger.Error().Msgf("Objective failed: %v", *m)
	})

	peers := []p2pms.PeerInfo{
		{Id: msgA.Id(), IpAddress: "127.0.0.1", Port: 3005, Address: alice.Address()},
		{Id: msgB.Id(), IpAddress: "127.0.0.1", Port: 3006, Address: bob.Address()},
	}

	// Connect nitro P2P message services
	msgA.AddPeers(peers)
	msgB.AddPeers(peers)

	defer msgA.Close()
	defer msgB.Close()

	trpA := transport.NewNatsTransport(nc, []string{
		fmt.Sprintf("nitro.%s", rpcproto.DirectFundRequestMethod),
		fmt.Sprintf("nitro.%s", rpcproto.DirectDefundRequestMethod),
	})
	conA, err := trpA.PollConnection()
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	ntsA := network.NewNetworkService(conA, &serde.MsgPack{})

	ntsA.RegisterRequestHandler(rpcproto.DirectFundRequestMethod, func(m *netproto.Message) {
		defer wg.Done()
		if len(m.Args) < 1 {
			logger.Fatal().Msg("unexpected empty args for direct funding method")
			return
		}

		for i := 0; i < len(m.Args); i++ {
			res := m.Args[i].(map[string]interface{})
			req := rpcproto.CreateDirectFundObjectiveRequest(res)

			assert.Equal(t, req.CounterParty, alice.Address())
			assert.Equal(t, req.AppDefinition, chainServiceA.GetConsensusAppAddress())

			nts.SendMessage(rpcproto.CreateDirectFundResponse(m.RequestId, &directfund.ObjectiveResponse{
				Id:        protocols.ObjectiveId("test"),
				ChannelId: types.Destination{},
			}))
			clientA.Engine.ObjectiveRequestsFromAPI <- req
		}
	})

	ntsA.RegisterRequestHandler(rpcproto.DirectDefundRequestMethod, func(m *netproto.Message) {
		defer wg.Done()
		if len(m.Args) < 1 {
			logger.Fatal().Msg("unexpected empty args for direct defunding method")
			return
		}

		for i := 0; i < len(m.Args); i++ {
			res := m.Args[i].(map[string]interface{})
			req := rpcproto.CreateDirectDefundObjectiveRequest(res)

			clientB.Engine.ObjectiveRequestsFromAPI <- req
		}
	})

	nts.SendMessage(
		rpcproto.CreateDirectFundRequestMessage(
			&directfund.ObjectiveRequest{
				CounterParty:      alice.Address(),
				ChallengeDuration: 100,
				Outcome:           testdata.Outcomes.Create(alice.Address(), bob.Address(), 100, 100, types.Address{}),
				AppDefinition:     chainServiceA.GetConsensusAppAddress(),
				// Appdata implicitly zero
				Nonce: rand.Uint64(),
			},
		),
	)
	// TODO: where do we create channel ids?
	nts.SendMessage(rpcproto.CreateDirectDefundRequestMessage(&directdefund.ObjectiveRequest{ChannelId: types.Destination{}}))
	wg.Add(2)

	wg.Wait()
}
