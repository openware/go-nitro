package rpcproto

import (
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/protocols"
)

type ObjectiveFailed struct {
	ObjectiveId protocols.ObjectiveId
	Reason      string
}

var _ netproto.Message = (*ObjectiveFailed)(nil)

const ObjectiveFailedType = "objective_failed"

func (r *ObjectiveFailed) Type() string {
	return ObjectiveFailedType
}
