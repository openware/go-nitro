package rpcproto

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/protocols/directfund"
)

const DirectFundRequestMethod = "direct_fund"

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

func CreateDirectFundRequest(r *directfund.ObjectiveRequest) *DirectFundRequest {
	o := []SingleAssetExit{}

	for _, ae := range r.Outcome {
		allocations := []Allocation{}
		for _, a := range ae.Allocations {
			allocations = append(allocations, Allocation{
				Destination:    a.Destination.String(),
				Amount:         a.Amount.String(),
				AllocationType: uint8(a.AllocationType),
				Metadata:       a.Metadata,
			})
		}
		o = append(o, SingleAssetExit{
			Asset:       ae.Asset.String(),
			Metadata:    ae.Metadata,
			Allocations: allocations,
		})
	}

	return &DirectFundRequest{
		CounterParty:      r.CounterParty.String(),
		ChallengeDuration: r.ChallengeDuration,
		Outcome:           o,
		AppDefinition:     r.AppDefinition.String(),
		AppData:           r.AppData,
		Nonce:             r.Nonce,
	}
}

func CreateObjectiveRequest(m map[string]interface{}) *directfund.ObjectiveRequest {
	r := directfund.ObjectiveRequest{
		CounterParty:      common.HexToAddress(m["counter_party"].(string)),
		ChallengeDuration: I2Uint32(m["challenge_duration"]),
		// Outcome:           ,
		AppDefinition: common.HexToAddress(m["app_definition"].(string)),
		AppData:       m["app_data"].([]byte),
		Nonce:         I2Uint64(m["nonce"]),
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
