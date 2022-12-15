package parser

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

func createAllocations(allocationInterfaces []any) []outcome.Allocation {
	allocationsArray := make([]outcome.Allocation, len(allocationInterfaces))
	for i := 0; i < len(allocationInterfaces); i++ {
		alloc := allocationInterfaces[i].(map[string]any)
		allocationsArray[i] = outcome.Allocation{
			Destination:    types.AddressToDestination(common.HexToAddress(alloc["destination"].(string))),
			Amount:         i2Uint256(alloc["amount"]),
			AllocationType: outcome.AllocationType(i2Uint8(alloc["allocation_type"])),
			Metadata:       toByteArray(alloc["metadata"]),
		}
	}
	return allocationsArray
}

func createExit(outcomesInterfaces []any) outcome.Exit {
	var e = outcome.Exit{}
	for i := 0; i < len(outcomesInterfaces); i++ {
		out := outcomesInterfaces[i].(map[string]any)
		allocations := out["allocations"].([]any)
		allocationsArray := createAllocations(allocations)

		e = append(e, outcome.SingleAssetExit{
			Asset:       common.HexToAddress(out["asset"].(string)),
			Metadata:    toByteArray(out["metadata"]),
			Allocations: allocationsArray,
		})
	}

	return e
}

func ParseDirectFundRequest(data map[string]any) directfund.ObjectiveRequest {
	outcomes := data["outcome"].([]any)

	r := directfund.ObjectiveRequest{
		CounterParty:      common.HexToAddress(data["counter_party"].(string)),
		ChallengeDuration: i2Uint32(data["challenge_duration"]),
		Outcome:           createExit(outcomes),
		AppDefinition:     common.HexToAddress(data["app_definition"].(string)),
		AppData:           toByteArray(data["app_data"]),
		Nonce:             i2Uint64(data["nonce"]),
	}

	return r
}

func ParseDirectFundResponse(data map[string]any) directfund.ObjectiveResponse {
	r := directfund.ObjectiveResponse{
		Id:        protocols.ObjectiveId(data["id"].(string)),
		ChannelId: types.Destination(common.HexToHash(data["channel_id"].(string))),
	}

	return r
}

func ParseDirectDefundRequest(data map[string]any) directdefund.ObjectiveRequest {
	r := directdefund.ObjectiveRequest{
		ChannelId: types.Destination(common.HexToHash(data["channel_id"].(string))),
	}

	return r
}

func ParseVirtualFundRequest(data map[string]any) virtualfund.ObjectiveRequest {
	outcomes := data["outcome"].([]any)

	r := virtualfund.ObjectiveRequest{
		Intermediaries:    hexesToAddresses(data["intermediaries"].([]string)),
		CounterParty:      common.HexToAddress(data["intermediaries"].(string)),
		ChallengeDuration: i2Uint32(data["challenge_duration"]),
		Outcome:           createExit(outcomes),
		Nonce:             i2Uint64(data["nonce"]),
		AppDefinition:     common.HexToAddress(data["app_definition"].(string)),
	}

	return r
}

func ParseVirtualFundResponse(data map[string]any) virtualfund.ObjectiveResponse {
	r := virtualfund.ObjectiveResponse{
		Id:        protocols.ObjectiveId(data["id"].(string)),
		ChannelId: types.Destination(common.HexToHash(data["channel_id"].(string))),
	}

	return r
}

func ParseVirtualDefundRequest(data map[string]any) virtualdefund.ObjectiveRequest {
	r := virtualdefund.ObjectiveRequest{
		ChannelId: types.Destination(common.HexToHash(data["channel_id"].(string))),
	}

	return r
}
