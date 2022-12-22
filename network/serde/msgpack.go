package serde

import (
	netproto "github.com/statechannels/go-nitro/network/protocol"
)

type MsgPack struct{}

func (j *MsgPack) Serialize(m *netproto.Message) ([]byte, error) {
	return m.MarshalMsg(nil)
}

func (j *MsgPack) Deserialize(data []byte) (*netproto.Message, error) {
	m := netproto.Message{}
	_, err := m.UnmarshalMsg(data)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
