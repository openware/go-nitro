package serde

import (
	"testing"

	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgPackSerializeRequest(t *testing.T) {
	rpc := MsgPack{}
	m := &netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: 4242,
		Method:    "test",
		Args:      []interface{}{"foo"},
	}
	data, err := rpc.Serialize(m)
	require.NoError(t, err)
	m2, err := rpc.Deserialize(data)
	require.NoError(t, err)
	assert.Equal(t, m, m2)
}

func TestMsgPackSerializeResponse(t *testing.T) {
	rpc := MsgPack{}
	m := &netproto.Message{
		Type:      netproto.TypeResponse,
		RequestId: 4242,
		Args:      []interface{}{"foo"},
	}
	data, err := rpc.Serialize(m)
	require.NoError(t, err)
	m2, err := rpc.Deserialize(data)
	require.NoError(t, err)
	assert.Equal(t, m, m2)
}

func TestMsgPackSerializeError(t *testing.T) {
	rpc := MsgPack{}
	m := &netproto.Message{
		Type:      netproto.TypeError,
		RequestId: 123,
		Method:    "test",
		Args:      []interface{}{int64(-32601), "Method not found"},
	}
	data, err := rpc.Serialize(m)
	require.NoError(t, err)
	m2, err := rpc.Deserialize(data)
	require.NoError(t, err)
	assert.Equal(t, m, m2)
}

func BenchmarkMsgPackSerializeRequest(b *testing.B) {
	rpc := MsgPack{}
	m := &netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: 4242,
		Method:    "test",
		Args:      []interface{}{"foo"},
	}

	for i := 0; i < b.N; i++ {
		rpc.Serialize(m)
	}
}

func BenchmarkMsgPackDeserializeRequest(b *testing.B) {
	rpc := MsgPack{}
	m := &netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: 4242,
		Method:    "test",
		Args:      []interface{}{"foo"},
	}
	data, err := rpc.Serialize(m)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		rpc.Deserialize(data)
	}
}
