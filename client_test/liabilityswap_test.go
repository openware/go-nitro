package client_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	Token "github.com/statechannels/go-nitro/client/engine/chainservice/erc20"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/require"
)

var liabilitySwapChannelDeposit = big.NewInt(1000000000000000000) // 1 eth in wei

// FIXME:
// 1. Test fails during the ledger channel funding with this error:
// 		0x111A00868581f73AB42FEEF67D235Ca09ca1E8db, error in run loop: failed to estimate gas needed: execution reverted: ERC20: transfer amount exceeds balance
//
// TODO:
// 1. Introduce the ability to use a Ledger Channel with custom App (AppDefinition and AppData)
//   a. Implement direct-advance protocol objective
//   b. Implement Engine advanceChannel function
//   c. Update Engine to support pure state channels (decouple from consensus channels, storage, etc.)
//   d. Check that direct direct-fund/defund protocols support pure state channels
// 2. Implement a Serde for liabilities
//

func TestLiabilitySwap(t *testing.T) {
	// Setup logging
	logFile := "test_liabilityswap.log"
	truncateLog(logFile)
	logger := newLogWriter(logFile)

	// Creates a new SimulatedBackend with the supplied number of transacting accounts
	// Deploys the Nitro Adjudicator
	// Deploys an ERC20 Token (Token WBTC)
	sim, bindings, accounts, err := chainservice.SetupSimulatedBackend(3)
	require.NoError(t, err)

	// Chain services setup
	chainI, err := chainservice.NewSimulatedBackendChainService(sim, bindings, accounts[1], logger)
	require.NoError(t, err)

	chainJ, err := chainservice.NewSimulatedBackendChainService(sim, bindings, accounts[2], logger)
	require.NoError(t, err)

	// Clients setup
	broker := messageservice.NewBroker()
	clientI, storeI := setupClient(irene.PrivateKey, chainI, broker, logger, 0)
	clientJ, storeJ := setupClient(brian.PrivateKey, chainJ, broker, logger, 0)

	// Tokens setup
	// tokenWBTCAddress := bindings.Token.Address

	// Deploy USDT ERC20 Token
	_, _, tokenUSDT, err := Token.DeployToken(accounts[0], sim, accounts[0].From)
	require.NoError(t, err)
	sim.Commit()

	_, err = tokenUSDT.Transfer(accounts[0], irene.Address(), big.NewInt(1_000_000))
	require.NoError(t, err)

	_, err = tokenUSDT.Transfer(accounts[0], brian.Address(), big.NewInt(1_000_000))
	require.NoError(t, err)
	sim.Commit()

	b, err := tokenUSDT.BalanceOf(&bind.CallOpts{}, irene.Address())
	require.NoError(t, err)
	logger.WriteString(fmt.Sprintf("Irene USDT balance: %d\n", b.Uint64()))

	// Deploy WETH ERC20 Token
	// tokenWETHAddress, _, _, err := Token.DeployToken(accounts[0], sim, accounts[0].From)
	// require.NoError(t, err)

	cId := liabilitySwapDirectlyFundALedgerChannel(
		t,
		clientI,
		clientJ,
		chainI.GetConsensusAppAddress(), // TODO: Use liabilities app instead of consensus app
		common.HexToAddress("0x0"),      // ETH (Native token)
	)

	want := createLiabilitySwapOutcome(*clientI.Address, *clientJ.Address, common.HexToAddress("0x0"), liabilitySwapChannelDeposit, liabilitySwapChannelDeposit)

	// Ensure that we create a consensus channel in the store
	for _, store := range []store.Store{storeI, storeJ} {
		var con *consensus_channel.ConsensusChannel
		var ok bool

		// each client fetches the ConsensusChannel by reference to their counterparty
		if store.GetChannelSecretKey() == &irene.PrivateKey {
			con, ok = store.GetConsensusChannel(*clientJ.Address)
		} else {
			con, ok = store.GetConsensusChannel(*clientI.Address)
		}
		require.True(t, ok, "expected a consensus channel to have been created")

		vars := con.ConsensusVars()
		got := vars.Outcome.AsOutcome()

		require.Equal(t, want, got, "unexpected outcome")
		//require.Equal(t, cmp.Diff(want, got), "unexpected outcome")

		require.Equal(t, uint64(1), vars.TurnNum, "expected consensus turn number to be the post fund setup 1")
		require.Equal(t, *clientI.Address, con.Leader())
		require.True(t, con.OnChainFunding.IsNonZero(), "Expected nonzero on chain funding, but got zero")

		// FIXME: why do we expect the channel to be destroyed at this stage?
		// _, channelStillInStore := store.GetChannelById(con.Id)
		// require.True(t, channelStillInStore, "Expected channel to have been destroyed")
	}

	// TODO: Advance channel using outcomes/appData to represent a swap
	_ = cId
	// oId := clientI.AdvanceChannel(
	// 	cId,
	// 	nextOutcome,
	// 	nextAppData,
	// )
	// waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, oId)
	// waitTimeForCompletedObjectiveIds(t, &clientJ, defaultTimeout, oId)
}

func liabilitySwapDirectlyFundALedgerChannel(
	t *testing.T,
	clientI, clientJ client.Client,
	appDeifinition types.Address,
	tokenAddress common.Address,
) types.Destination {
	// Set up an outcome that requires both participants to deposit
	outcome := createLiabilitySwapOutcome(*clientI.Address, *clientJ.Address, tokenAddress, liabilitySwapChannelDeposit, liabilitySwapChannelDeposit)

	appData := []byte("") // TODO: Create initial appData once format is decided

	response := clientI.CreateCustomLedgerChannel(
		*clientJ.Address,
		appDeifinition,
		0,
		outcome,
		appData,
	)

	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, response.Id)
	waitTimeForCompletedObjectiveIds(t, &clientJ, defaultTimeout, response.Id)
	return response.ChannelId
}

// createLiabilitySwapOutcome is a helper function to create a two-actor outcome
func createLiabilitySwapOutcome(destA, destB, tokenAddress types.Address, amountA, amountB *big.Int) outcome.Exit {
	return outcome.Exit{
		outcome.SingleAssetExit{
			Asset: tokenAddress,
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: types.AddressToDestination(destA),
					Amount:      amountA,
				},
				outcome.Allocation{
					Destination: types.AddressToDestination(destB),
					Amount:      amountB,
				},
			},
		},
	}
}
