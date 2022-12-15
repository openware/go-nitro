package rpcproto

import (
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/protocols/directfund"
)

type DirectFundResponse struct {
	netproto.Response

	Result directfund.ObjectiveResponse
}

var _ netproto.ResponseMessage = (*DirectFundResponse)(nil)

const DirectFundResponseType = "direct_fund_response"

func (r *DirectFundResponse) Type() string {
	return DirectFundResponseType
}
