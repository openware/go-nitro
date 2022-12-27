package rpcproto

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

//go:generate msgp

type SingleAssetExit struct {
	Asset       string       `msg:"asset"`    // Either the zero address (implying the native token) or the address of an ERC20 contract
	Metadata    []byte       `msg:"metadata"` // Can be used to encode arbitrary additional information that applies to all allocations.
	Allocations []Allocation `msg:"allocations"`
}

type Allocation struct {
	Destination    string `msg:"destination"`     // Either an ethereum address or an application-specific identifier
	Amount         string `msg:"amount"`          // An amount of a particular asset
	AllocationType uint8  `msg:"allocation_type"` // Directs calling code on how to interpret the allocation
	Metadata       []byte `msg:"metadata"`        // Custom metadata (optional field, can be zero bytes). This can be used flexibly by different protocols.
}

func createAllocations(allocationInterfaces []interface{}) []outcome.Allocation {
	allocationsArray := make([]outcome.Allocation, len(allocationInterfaces))
	for i := 0; i < len(allocationInterfaces); i++ {
		alloc := allocationInterfaces[0].(map[string]interface{})
		allocationsArray[i] = outcome.Allocation{
			Destination:    types.AddressToDestination(common.HexToAddress(alloc["destination"].(string))),
			Amount:         I2Uint256(alloc["amount"].(string)),
			AllocationType: outcome.AllocationType(I2Uint8(alloc["allocation_type"])),
			Metadata:       alloc["metadata"].([]byte),
		}
	}
	return allocationsArray
}

func createExit(outcomesInterfaces []interface{}) outcome.Exit {
	var e = outcome.Exit{}
	for i := 0; i < len(outcomesInterfaces); i++ {
		out := outcomesInterfaces[0].(map[string]interface{})
		allocations := out["allocations"].([]interface{})
		allocationsArray := createAllocations(allocations)

		e = append(e, outcome.SingleAssetExit{
			Asset:       common.HexToAddress(out["asset"].(string)),
			Metadata:    out["metadata"].([]byte),
			Allocations: allocationsArray,
		})
	}

	return e
}
