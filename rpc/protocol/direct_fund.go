package rpcproto

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/protocols/directfund"
)

// TODO: maybe 1 for request and 1 for response
const DirectFundRequestMethod = "direct_fund"

//go:generate msgp

type DirectFundRequest struct {
	CounterParty      string            `msg:"counter_party"`
	ChallengeDuration uint32            `msg:"challenge_duration"`
	Outcome           []SingleAssetExit `msg:"outcome"`
	AppDefinition     string            `msg:"app_definition"`
	AppData           []byte            `msg:"app_data"`
	Nonce             uint64            `msg:"nonce"`
}

type DirectFundResponse struct {
	Id        string `msg:"id"`
	ChannelId string `msg:"channel_id"`
}

func CreateDirectFundObjectiveRequest(m map[string]interface{}) *directfund.ObjectiveRequest {
	outcomes := m["outcome"].([]interface{})
	exit := createExit(outcomes)

	r := directfund.ObjectiveRequest{
		CounterParty:      common.HexToAddress(m["counter_party"].(string)),
		ChallengeDuration: I2Uint32(m["challenge_duration"]),
		Outcome:           exit,
		AppDefinition:     common.HexToAddress(m["app_definition"].(string)),
		AppData:           m["app_data"].([]byte),
		Nonce:             I2Uint64(m["nonce"]),
	}

	return &r
}

func CreateDirectFundRequestMessage(r *directfund.ObjectiveRequest) *netproto.Message {
	return &netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: rand.Uint64(),
		Method:    DirectFundRequestMethod,
		Args:      []interface{}{CreateDirectFundRequest(r)},
	}
}

func CreateDirectFundResponse(reqId uint64, args *directfund.ObjectiveResponse) *netproto.Message {
	r := DirectFundResponse{
		Id:        string(args.Id),
		ChannelId: args.ChannelId.String(),
	}
	return &netproto.Message{
		Type:      netproto.TypeResponse,
		RequestId: reqId,
		Method:    DirectFundRequestMethod,
		Args:      []interface{}{&r},
	}
}
