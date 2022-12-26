package rpcproto

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/types"
)

const DirectDefundRequestMethod = "direct_defund"

//go:generate msgp

type DirectDefundRequest struct {
	ChannelId string `msg:"channel_id"`
}

func CreateDirectDefundObjectiveRequest(m map[string]interface{}) *directdefund.ObjectiveRequest {
	r := directdefund.ObjectiveRequest{
		ChannelId: types.AddressToDestination(common.HexToAddress(m["channel_id"].(string))),
	}

	return &r
}

func CreateDirectDefundRequestMessage(args *directdefund.ObjectiveRequest) *netproto.Message {
	r := DirectDefundRequest{
		ChannelId: args.ChannelId.String(),
	}

	return &netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: rand.Uint64(),
		Method:    DirectDefundRequestMethod,
		Args:      []interface{}{&r},
	}
}
