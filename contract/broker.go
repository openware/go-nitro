package contract

import (
	"github.com/statechannels/go-nitro/types"
)

type Broker struct {
	Address     types.Address
	Destination types.Destination
	PrivateKey  []byte
	Role        uint
}

func NewBroker(address types.Address, destination types.Destination, privateKey []byte, role uint) (*Broker, error) {
	return &Broker{Address: address, Destination: destination, PrivateKey: privateKey, Role: role}, nil
}
