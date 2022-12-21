package rpcproto

import (
	"math/rand"

	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/protocols/directfund"
)

const DirectFundRequestMethod = "direct_fund"

func CreateDirectFundRequest(args *directfund.ObjectiveRequest) *netproto.Message {
	return &netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: rand.Uint64(),
		Method:    DirectFundRequestMethod,
		Args:      []interface{}{*args},
	}
}

func CreateDirectFundResponse(reqId uint64, args *directfund.ObjectiveResponse) *netproto.Message {
	return &netproto.Message{
		Type:      netproto.TypeResponse,
		RequestId: reqId,
		Method:    DirectFundRequestMethod,
		Args:      []interface{}{*args},
	}
}
