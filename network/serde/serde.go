package serde

import netproto "github.com/statechannels/go-nitro/network/protocol"

type Serde interface {
	Serializer
	Deserializer
}

type Serializer interface {
	Serialize(*netproto.Message) ([]byte, error)
}

type Deserializer interface {
	Deserialize([]byte) (*netproto.Message, error)
}

// JSON RPC example:
//
// payload = {
// 	"method": method,
// 	"params": [args],
// 	"jsonrpc": "2.0",
// 	"id": 1,
// }

// JSON array RPC:
//
// [type, request_id, method, arguments]
