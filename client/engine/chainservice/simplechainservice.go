package chainservice

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// SimpleChainService forwards inputted transactions to a MockChain, and passes Events straight back.
type SimpleChainService struct {
	out chan Event                 // out is the chan used to send Events to the engine
	in  chan protocols.Transaction // in is the chan used to recieve Transactions from the engine

	address types.Address // address is used to subscribe to the MockChain's Out chan
	chain   MockChain
}

// NewSimpleChainService returns a SimpleChainService which is listening for transactions and events.
func NewSimpleChainService(mc MockChain, address types.Address) ChainService {
	mcs := SimpleChainService{}
	mcs.out = make(chan Event)
	mcs.in = make(chan protocols.Transaction)
	mcs.chain = mc
	mcs.address = address

	go mcs.listenForEvents()
	go mcs.listenForTransactions()

	return mcs
}

// Out() returns the but chan but narrows the type so that external consumers mays only recieve on it.
func (mcs SimpleChainService) Out() <-chan Event {
	return mcs.out
}

// In returns the in chan but narrows the type so that external consumers mays only send on it.
func (mcs SimpleChainService) In() chan<- protocols.Transaction {
	return mcs.in
}

// listenForTransactions pipes transactions to the MockChain
func (mcs SimpleChainService) listenForTransactions() {
	for tx := range mcs.in {
		mcs.chain.In() <- tx
	}
}

// listenForEvents peipes events from the MockChain
func (mcs SimpleChainService) listenForEvents() {
	for event := range mcs.chain.Out(mcs.address) {
		mcs.out <- event
	}
}
