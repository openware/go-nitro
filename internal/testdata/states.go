package testdata

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

// a simpleItem is an ergonomic way to create a simple allocation item, without having to create a big.Int
// This is more ergonomic for creating test data
type simpleItem struct {
	Dest   types.Destination
	Amount int64
}

type outcomes struct {
	// Create returns a simple outcome {a: aBalance, b: bBalance} in the
	// zero-asset (chain-native token)
	Create func(a, b types.Address, aBalance, bBalance uint) outcome.Exit
	// CreateLongOutcome returns a simple outcome {addressOne: balanceOne ...} in the
	// zero-asset (chain-native token)
	CreateLongOutcome func(items ...simpleItem) outcome.Exit
}

var Outcomes outcomes = outcomes{
	Create:            createOutcome,
	CreateLongOutcome: createLongOutcome,
}

var chainId, _ = big.NewInt(0).SetString("9001", 10)
var someAppDefinition = common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`)

var testOutcome = outcome.Exit{
	outcome.SingleAssetExit{
		Asset: types.Address{}, // the native token of the chain
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: Actors.Alice.Destination(),
				Amount:      big.NewInt(6),
			},
			outcome.Allocation{
				Destination: Actors.Bob.Destination(),
				Amount:      big.NewInt(4),
			},
		},
	},
}

var testVirtualState = state.State{
	ChainId: chainId,
	Participants: []types.Address{
		Actors.Alice.Address,
		Actors.Irene.Address,
		Actors.Bob.Address,
	},
	ChannelNonce:      big.NewInt(1234789),
	AppDefinition:     someAppDefinition,
	ChallengeDuration: big.NewInt(60),
	AppData:           []byte{},
	Outcome: Outcomes.CreateLongOutcome(
		simpleItem{Actors.Alice.Destination(), 6},
		simpleItem{Actors.Bob.Destination(), 4},
	),
	TurnNum: 0,
	IsFinal: false,
}

var testState = state.State{
	ChainId: chainId,
	Participants: []types.Address{
		Actors.Alice.Address,
		Actors.Bob.Address,
	},
	ChannelNonce:      big.NewInt(37140676580),
	AppDefinition:     someAppDefinition,
	ChallengeDuration: big.NewInt(60),
	AppData:           []byte{},
	Outcome:           testOutcome,
	TurnNum:           5,
	IsFinal:           false,
}

func createLedgerState(client, hub types.Address, clientBalance, hubBalance uint) state.State {
	state := testState.Clone()
	state.Participants = []types.Address{
		client,
		hub,
	}
	state.Outcome = Outcomes.Create(client, hub, clientBalance, hubBalance)
	state.AppDefinition = types.Address{} // ledger channel running the consensus app
	state.TurnNum = 0                     // a requirement for channel.NewTwoPartyLedger

	return state
}

// createOutcome is a helper function to create a two-actor outcome
func createOutcome(first types.Address, second types.Address, x, y uint) outcome.Exit {

	return outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: types.AddressToDestination(first),
				Amount:      big.NewInt(int64(x)),
			},
			outcome.Allocation{
				Destination: types.AddressToDestination(second),
				Amount:      big.NewInt(int64(y)),
			},
		},
	}}
}

func createLongOutcome(items ...simpleItem) outcome.Exit {
	sae := outcome.SingleAssetExit{}
	for _, i := range items {
		a := outcome.Allocation{
			Destination: i.Dest,
			Amount:      big.NewInt(i.Amount),
		}
		sae.Allocations = append(sae.Allocations, a)
	}

	return outcome.Exit{sae}
}
