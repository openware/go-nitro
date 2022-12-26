package rpcproto

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/types"
)

const VirtualDefundRequestMethod = "virtual_defund"

//go:generate msgp

type VirtualDefundRequest struct {
	ChannelId string `msg:"channel_id"`
}

func CreateVirtualDefundObjectiveRequest(m map[string]interface{}) *virtualdefund.ObjectiveRequest {
	r := virtualdefund.ObjectiveRequest{
		ChannelId: types.AddressToDestination(common.HexToAddress(m["channel_id"].(string))),
	}

	return &r
}

func CreateVirtualDefundRequestMessage(args *virtualdefund.ObjectiveRequest) *netproto.Message {
	r := VirtualDefundRequest{
		ChannelId: args.ChannelId.String(),
	}

	return &netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: rand.Uint64(),
		Method:    VirtualDefundRequestMethod,
		Args:      []interface{}{&r},
	}
}
