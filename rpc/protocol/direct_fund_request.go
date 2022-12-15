package rpcproto

import (
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/protocols/directfund"
)

type DirectFundRequest struct {
	netproto.Request

	Params directfund.ObjectiveRequest
}

var _ netproto.RequestMessage = (*DirectFundRequest)(nil)

const DirectFundRequestType = "direct_fund_request"

func (r *DirectFundRequest) Type() string {
	return DirectFundRequestType
}
