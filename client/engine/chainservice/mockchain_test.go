package chainservice

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestDeposit(t *testing.T) {
	// The MockChain and SimpleChainService should work together to react to a deposit transaction for a given channel by:
	// - sending an event with updated holdings for that channel to all SimpleChainServices which are subscribed

	var a = types.Address(common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`))
	var b = types.Address(common.HexToAddress(`0xa5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`))

	// Construct MockChain and tell it the addresses of the SimpleChainServices which will subscribe to it.
	// This is not super elegant but gets around data races -- the constructor will make channels and then run a listener which will send on them.
	var chain = NewMockChain([]types.Address{a, b})

	// Construct SimpleChainServices
	mcsA := NewSimpleChainService(chain, a)
	mcsB := NewSimpleChainService(chain, b)

	inA := mcsA.In()
	outA := mcsA.Out()

	// Prepare test data to trigger MockChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(1),
	}
	testTx := protocols.Transaction{
		ChannelId: types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`)),
		Deposit:   testDeposit,
	}

	// Send one transaction into one of the SimpleChainServices and recieve one event from it.
	inA <- testTx
	event := <-outA

	if event.ChannelId != testTx.ChannelId {
		t.Error(`channelId mismatch`)
	}
	if !event.Holdings.Equal(testTx.Deposit) {
		t.Error(`holdings mismatch`)
	}

	// Send the transaction again and recieve another event
	inA <- testTx
	event = <-outA

	// The expectation is that the MockChainService remembered the previous deposit and added this one to it:
	expectedHoldings := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(2),
	}

	if event.ChannelId != testTx.ChannelId {
		t.Error(`channelId mismatch`)
	}
	if !event.Holdings.Equal(expectedHoldings) {
		t.Error(`holdings mismatch`)
	}

	// Pull an event out of the other mock chain service and check that
	eventB := <-mcsB.Out()

	if eventB.ChannelId != testTx.ChannelId {
		t.Error(`channelId mismatch`)
	}
	if !eventB.Holdings.Equal(testTx.Deposit) {
		t.Error(`holdings mismatch`)
	}

	// Pull another event out of the other mock chain service and check that
	eventB = <-mcsB.Out()

	if eventB.ChannelId != testTx.ChannelId {
		t.Error(`channelId mismatch`)
	}
	if !eventB.Holdings.Equal(expectedHoldings) {
		t.Error(`holdings mismatch`)
	}

}
