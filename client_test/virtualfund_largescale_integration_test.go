package client_test

import (
	"encoding/json"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	nc "github.com/statechannels/go-nitro/crypto"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

func TestLargeScaleVirtualFundIntegration(t *testing.T) {
	const numRetrievalClients = 1

	logDir := "../artifacts/vectorclock"
	logFile := "largescale_client_test.log"

	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	retrievalProvider, retrievalProviderStore := setupInstrumentedClient(bob.PrivateKey, chain, broker, logDestination, 0, logDir)
	paymentHub, _ := setupInstrumentedClient(irene.PrivateKey, chain, broker, logDestination, 0, logDir)

	directlyFundALedgerChannel(t, retrievalProvider, paymentHub)

	retrievalClients := make([]client.Client, numRetrievalClients)
	for i, _ := range retrievalClients {
		secretKey, _ := nc.GeneratePrivateKeyAndAddress()
		retrievalClients[i], _ = setupInstrumentedClient(secretKey, chain, broker, logDestination, 0, logDir)
		directlyFundALedgerChannel(t, retrievalClients[i], paymentHub)
		go createVirtualChannelWithRetrievalProvider(retrievalClients[i], retrievalProvider)
	}

	<-time.After(5 * time.Second)

	retrievalProviderHubConnection, _ := retrievalProviderStore.GetConsensusChannel(*paymentHub.Address)

	finalOutcome, _ := json.Marshal(retrievalProviderHubConnection.SupportedSignedState().State().Outcome)

	logDestination.Write(finalOutcome)

}

func createVirtualChannelWithRetrievalProvider(c client.Client, retrievalProvider client.Client) protocols.ObjectiveId {
	withRetrievalProvider := virtualfund.ObjectiveRequest{
		MyAddress:    *c.Address,
		CounterParty: *retrievalProvider.Address,
		Intermediary: irene.Address(),
		Outcome: td.Outcomes.Create(
			*c.Address,
			*retrievalProvider.Address,
			1,
			1,
		),
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}
	return c.CreateVirtualChannel(withRetrievalProvider)
}
