package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/contract"
	typ "github.com/statechannels/go-nitro/types"
)

var (
	DefaultTimeout = 3 * time.Second
	ChainId        = big.NewInt(1)
	RpcUrl         = "http://127.0.0.1:8545"

	// TODO
	broker1, _ = contract.NewBroker(
		common.HexToAddress(`0x7d6fe92F348B6F2216A7AA2c2F0Dd9b8c830e490`),
		typ.AddressToDestination(common.HexToAddress(`0x7d6fe92F348B6F2216A7AA2c2F0Dd9b8c830e490`)),
		common.Hex2Bytes(`9a14bf0eb618a3407a12a83a74dfe7bbed098ccc6347985b92ab08e81996cfc9`),
		0,
	)

	broker2, _ = contract.NewBroker(
		common.HexToAddress(`0xb1239c28162bf9b3e2aa6Dc2c78066B26D5423F7`),
		typ.AddressToDestination(common.HexToAddress(`0xb1239c28162bf9b3e2aa6Dc2c78066B26D5423F7`)),
		common.Hex2Bytes(`2305f34d1dcab90a5143856446d5213c0b29ae353a25845445c8050d4bca38d9`),
		1,
	)
)

func sign(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	signer := types.LatestSignerForChainID(ChainId)
	hash := signer.Hash(tx)

	var privateKey []byte
	// TODO
	if address == broker1.Address {
		privateKey = broker1.PrivateKey
	} else {
		privateKey = broker2.PrivateKey
	}

	prv, err := crypto.ToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	signature, err := crypto.Sign(hash.Bytes(), prv)
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(signer, signature)
}

//  Broker A create initial channel state
//  Broker A deposit funds to state channel using channel id generated from initial state
//  Broker B get initial channel state from Broker A
//  Broker B check if Broker A deposited specified amount in outcome, if so continue
//  Broker B deposit to state channel
//  Broker B update state with his new outcomes
//  Broker A reviece state update
//  Broker A check if Broker B outcome is the same as what is inside state channel smart contract
//  If so, broker A agree to continue working on this state channel
func main() {
	// STEP 1 - deploy smart contract (nitro adjucator)

	// STEP 2 - initialize client
	client, err := contract.NewClient("0xCc388ae2496E15ff8C6df70566171c750B5118E2", "http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}

	// STEP 3 - Create brokers
	// There are 2 participants (broker 1 and broker 2) in the system

	// STEP 4 - Open a channel between participants
	chainId, err := client.Contract.GetChainID(nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("chainId: %v\n\n", chainId)

	// initial state to sign
	var preFundState = state.State{
		ChainId:           chainId,
		Participants:      []typ.Address{broker1.Address, broker2.Address},
		ChannelNonce:      big.NewInt(0), //big.NewInt(time.Now().UnixNano()),
		ChallengeDuration: big.NewInt(60),
		AppData:           []byte{},
		Outcome: outcome.Exit{
			outcome.SingleAssetExit{
				Asset: common.HexToAddress("0x00"),
				Allocations: outcome.Allocations{
					outcome.Allocation{
						Destination: broker1.Destination,
						Amount:      big.NewInt(1),
					},
					outcome.Allocation{
						Destination: broker2.Destination,
						Amount:      big.NewInt(2),
					},
				},
			},
		},
		TurnNum: 0,
		IsFinal: false,
	}

	c, err := channel.New(preFundState, broker1.Role)
	if err != nil {
		panic(err)
	}

	// STEP 5 - Sign prefund state
	signature1, err := preFundState.Sign(broker1.PrivateKey)
	c.AddSignedState(preFundState, signature1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Signature 1 %+v\n\n", signature1)

	signature2, err := preFundState.Sign(broker2.PrivateKey)
	c.AddSignedState(preFundState, signature2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Signature 2 %+v\n\n", signature2)
	fmt.Printf("Channel ID: %+v \n\n", c.Id)

	// STEP 6 - Deposit process
	// Deposit funds only of prefund state complete
	if c.PreFundComplete() {
		signerFn := sign
		transactionOpts1 := bind.TransactOpts{From: broker1.Address, Signer: signerFn, Value: big.NewInt(1)}
		transaction1, err := client.Contract.Deposit(&transactionOpts1, common.HexToAddress("0x00"), c.Id, big.NewInt(0), big.NewInt(1))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Transaction 1:  %v\n\n", transaction1)

		transactionOpts2 := bind.TransactOpts{From: broker2.Address, Signer: signerFn, Value: big.NewInt(2)}
		transaction2, err := client.Contract.Deposit(&transactionOpts2, common.HexToAddress("0x00"), c.Id, big.NewInt(1), big.NewInt(2))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Transaction 2:  %v\n\n", transaction2)
	}

	// STEP 7 - Sign post fund state
	signature1, err = c.PostFundState().Sign(broker1.PrivateKey)
	c.AddSignedState(c.PostFundState(), signature1)
	if err != nil {
		panic(err)
	}

	signature2, err = c.PostFundState().Sign(broker2.PrivateKey)
	c.AddSignedState(c.PostFundState(), signature2)
	if err != nil {
		panic(err)
	}

	if c.PostFundComplete() {
		// here they send info to each other
		var secondState = state.State{
			ChainId:           chainId,
			Participants:      []typ.Address{broker1.Address, broker2.Address},
			ChannelNonce:      big.NewInt(0),
			ChallengeDuration: big.NewInt(60),
			AppData:           []byte{}, // here additional info
			Outcome: outcome.Exit{
				outcome.SingleAssetExit{
					Asset: common.HexToAddress("0x00"),
					Allocations: outcome.Allocations{
						outcome.Allocation{
							Destination: broker1.Destination,
							Amount:      big.NewInt(2),
						},
						outcome.Allocation{
							Destination: broker2.Destination,
							Amount:      big.NewInt(1),
						},
					},
				},
			},
			TurnNum: 2,
			IsFinal: false,
		}

		signature1, err = c.CurrentState(secondState).Sign(broker1.PrivateKey)
		c.AddSignedState(c.CurrentState(secondState), signature1)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Current state completion after 1 participant: %v\n\n", c.CurrentStateComplete(secondState))
		fmt.Printf("State signed by broker1: %v\n\n", c.CurrentStateSignedByMe(secondState))

		signature2, err = c.CurrentState(secondState).Sign(broker2.PrivateKey)
		c.AddSignedState(c.CurrentState(secondState), signature2)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Current state completion after 2 participant: %v\n\n", c.CurrentStateComplete(secondState))

		fmt.Printf("Channel signatures: %+v\n\n", c.SignedStateForTurnNum)
	}
}
