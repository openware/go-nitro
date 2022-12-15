package rpcproto

import (
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/protocols"
)

type ObjectiveUpdated struct {
	ObjectiveId protocols.ObjectiveId
	Status      protocols.ObjectiveStatus
}

var _ netproto.Message = (*ObjectiveUpdated)(nil)

const ObjectiveUpdatedType = "objective_updated"

func (r *ObjectiveUpdated) Type() string {
	return ObjectiveUpdatedType
}
