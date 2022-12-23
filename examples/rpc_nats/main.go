package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

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
	"github.com/statechannels/go-nitro/protocols/directfund"
	rpcproto "github.com/statechannels/go-nitro/rpc/protocol"
	"github.com/statechannels/go-nitro/types"
)

var natsConnectionUrl = "localhost:4222"
var nc *nats.Conn

var (
	alice = testactors.Alice
	bob   = testactors.Bob
	irene = testactors.Irene
)

func init() {
	var err error
	nc, err = transport.InitNats(natsConnectionUrl)
	if err != nil {
		panic("can't connect to nats.")
	}
}

func setupClientWithP2PMessageService(pk []byte, port int, chain *chainservice.MockChainService, logDestination io.Writer) (client.Client, *p2pms.P2PMessageService) {
	messageservice := p2pms.NewMessageService("127.0.0.1", port, pk)
	storeA := store.NewMemStore(pk)
	return client.New(messageservice, chain, storeA, logDestination, &engine.PermissivePolicy{}, nil), messageservice
}

var (
	chain = chainservice.NewMockChain()

	chainServiceA = chainservice.NewMockChainService(chain, alice.Address())

	//trpA = transport.NewChanTransport()
)

// Go nitro micro service entry code example
// NOTE: this example is not accurate, since we have to add B and I clients in order to have a working example
// On actual go-nitro service, only "A" related code would be present
func nitroService(logger zerolog.Logger) {
	// Setup logger
	logger = logger.With().
		Str("client", "NITRO ").
		Str("scope", "     ").
		Logger()

	// logFile := "rpc_nats.log"
	// truncateLog(logFile)
	// logDestination := newLogWriter(logFile)

	// TODO: refactor rpc service to allow chain and P2P MS updates
	// for exemple: disconnect from B or I, reconnect to B or I, ...
	// Orverall, the goal is to be able to completly control the client trough the rpc service

	// Setup B and I clients
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())

	clientA, msgA := setupClientWithP2PMessageService(alice.PrivateKey, 3005, chainServiceA, logger)
	clientB, msgB := setupClientWithP2PMessageService(bob.PrivateKey, 3006, chainServiceB, logger)
	clientI, msgI := setupClientWithP2PMessageService(irene.PrivateKey, 3007, chainServiceI, logger)
	peers := []p2pms.PeerInfo{
		{Id: msgA.Id(), IpAddress: "127.0.0.1", Port: 3005, Address: alice.Address()},
		{Id: msgB.Id(), IpAddress: "127.0.0.1", Port: 3006, Address: bob.Address()},
		{Id: msgI.Id(), IpAddress: "127.0.0.1", Port: 3007, Address: irene.Address()},
	}

	// Connect nitro P2P message services
	msgA.AddPeers(peers)
	msgB.AddPeers(peers)
	msgI.AddPeers(peers)

	defer msgA.Close()
	defer msgB.Close()
	defer msgI.Close()

	// Ignore B and I clients for now
	_ = clientB
	_ = clientI

	// Setup A network service using transport from global variables (in global only because we currently use a mock transport)
	trpA := transport.NewNatsTransport(nc, []string{fmt.Sprintf("nitro.%s", rpcproto.DirectFundRequestMethod), "nitro.test-topic"})
	conA, err := trpA.PollConnection()
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	ntsA := network.NewNetworkService(conA, &serde.MsgPack{})
	ntsA.Logger = logger.With().Str("scope", "NETW ").Logger()

	ntsA.RegisterRequestHandler(rpcproto.DirectFundRequestMethod, func(m *netproto.Message) {
		if len(m.Args) < 1 {
			logger.Fatal().Msg("unexpected empty args for direct funding method")
			return
		}

		for i := 0; i < len(m.Args); i++ {
			res := m.Args[i].(map[string]interface{})
			req := rpcproto.CreateObjectiveRequest(res)

			logger.Info().Msgf("Objective Request: %v", req)
			clientA.Engine.ObjectiveRequestsFromAPI <- req
		}
	})

	// TODO: complete example with B and I clients interactions (wait their own objectives, etc.)
	ntsA.RegisterRequestHandler(rpcproto.DirectFundRequestMethod, func(m *netproto.Message) {
		if len(m.Args) < 1 {
			logger.Fatal().Msg("unexpected empty args for direct funding method")
			return
		}

		for i := 0; i < len(m.Args); i++ {
			res := m.Args[i].(map[string]interface{})
			req := rpcproto.CreateObjectiveRequest(res)

			logger.Info().Msgf("Objective Request: %v", req)
			clientB.Engine.ObjectiveRequestsFromAPI <- req
		}
	})

	// Wait forever
	select {}
}

// Simulated external micro service example
func marginService(logger zerolog.Logger) {
	// Setup logger
	logger = logger.With().
		Str("client", "MARGIN").
		Str("scope", "     ").
		Logger()

	// Setup transport
	//trp := transport.NewChanTransport()
	//trp.Connect(trpA, 100*time.Millisecond, 1000*time.Millisecond)
	trp := transport.NewNatsTransport(nc, []string{fmt.Sprintf("nitro.%s", rpcproto.DirectFundRequestMethod), "nitro.test-topic"})

	// Setup network service
	con, err := trp.PollConnection()
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	nts := network.NewNetworkService(con, &serde.MsgPack{})
	nts.Logger = logger.With().Str("scope", "NETW ").Logger()
	// TODO: if we close it's an error
	//defer nts.Close()

	// NOTE: instead of manually using network service, like bellow example, we could use the rpc service
	// instead, that will add helper methods with the same behavior
	// This would require external micro services to have a dependency on the rpc service, which may not be desirable

	nts.RegisterResponseHandler(rpcproto.DirectFundRequestMethod, func(m *netproto.Message) {
		logger.Info().Msgf("Objective updated: %v", *m)
	})

	nts.RegisterErrorHandler(rpcproto.DirectFundRequestMethod, func(m *netproto.Message) {
		logger.Error().Msgf("Objective failed: %v", *m)
	})

	// Start a new goroutine to handle the peer
	// Register objective failed handler

	// Send direct fund request
	nts.SendMessage(
		rpcproto.CreateDirectFundRequestMessage(
			&directfund.ObjectiveRequest{
				CounterParty:      irene.Address(),
				ChallengeDuration: 0,
				Outcome:           testdata.Outcomes.Create(alice.Address(), irene.Address(), 100, 100, types.Address{}),
				AppDefinition:     chainServiceA.GetConsensusAppAddress(),
				// Appdata implicitly zero
				Nonce: rand.Uint64(),
			},
		),
	)

	nts.SendMessage(
		rpcproto.CreateDirectFundRequestMessage(
			&directfund.ObjectiveRequest{
				CounterParty:      bob.Address(),
				ChallengeDuration: 100,
				Outcome:           testdata.Outcomes.Create(alice.Address(), bob.Address(), 100, 100, types.Address{}),
				AppDefinition:     chainServiceA.GetConsensusAppAddress(),
				// Appdata implicitly zero
				Nonce: rand.Uint64(),
			},
		),
	)
}

func main() {
	// Setup logger
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:           os.Stdout,
		TimeFormat:    time.RFC3339,
		PartsOrder:    []string{"time", "level", "caller", "client", "scope", "message"},
		FieldsExclude: []string{"time", "level", "caller", "message", "client", "scope"},
	}).
		// Level(zerolog.DebugLevel).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("client", "").
		Str("scope", "").
		Logger()

	// Start nitro micro service
	go nitroService(logger)

	// Start margin micro service (simulated external micro service)
	go func() {
		time.Sleep(time.Millisecond * 100)
		marginService(logger)
	}()

	// Wait forever
	select {}
}
