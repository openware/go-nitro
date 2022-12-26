package rpcproto

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

const VirtualFundRequestMethod = "virtual_fund"

//go:generate msgp

type VirtualFundRequest struct {
	Intermediaries    []string          `msg:"intermediaries"`
	CounterParty      string            `msg:"counter_party"`
	ChallengeDuration uint32            `msg:"challenge_duration"`
	Outcome           []SingleAssetExit `msg:"outcome"`
	Nonce             uint64            `msg:"nonce"`
	AppDefinition     string            `msg:"app_definition"`
}

type VirtualFundResponse struct {
	Id        string `msg:"id"`
	ChannelId string `msg:"channel_id"`
}

func CreateVirtualFundRequest(r *virtualfund.ObjectiveRequest) *VirtualFundRequest {

}

func CreateVirtualFundObjectiveRequest(m map[string]interface{}) *virtualfund.ObjectiveRequest {
	outcomes := m["outcome"].([]interface{})
	exit := createExit(outcomes)

	// TODO: maybe make a helper method
	intermediaries := m["counter_party"].([]string)
	intermediariesAddresses := make([]types.Address, len(intermediaries))
	for i := 0; i < len(intermediariesAddresses); i++ {
		intermediariesAddresses[i] = common.HexToAddress(intermediaries[i])
	}

	r := virtualfund.ObjectiveRequest{
		Intermediaries:    intermediariesAddresses,
		CounterParty:      common.HexToAddress(m["counter_party"].(string)),
		ChallengeDuration: I2Uint32(m["challenge_duration"]),
		Outcome:           exit,
		Nonce:             I2Uint64(m["nonce"]),
		AppDefinition:     common.HexToAddress(m["app_definition"].(string)),
	}

	return &r
}

func CreateVirtualFundRequestMessage(r *virtualfund.ObjectiveRequest) *netproto.Message {
	return &netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: rand.Uint64(),
		Method:    VirtualFundRequestMethod,
		Args:      []interface{}{CreateVirtualFundRequest(r)},
	}
}

func CreateVirtualFundResponseMessage(reqId uint64, args *virtualfund.ObjectiveResponse) *netproto.Message {
	r := VirtualFundResponse{
		Id:        string(args.Id),
		ChannelId: args.ChannelId.String(),
	}

	return &netproto.Message{
		Type:      netproto.TypeResponse,
		RequestId: reqId,
		Method:    VirtualFundRequestMethod,
		Args:      []interface{}{&r},
	}
}
