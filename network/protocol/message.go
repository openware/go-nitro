package netproto

import "strconv"

type MessageType int8

const (
	TypeRequest      MessageType = 1
	TypeResponse                 = 2
	TypePublicEvent              = 3
	TypePrivateEvent             = 4
	TypeError                    = 5
)

//go:generate msgp
type Message struct {
	Type      MessageType   `msg:"type"`
	RequestId uint64        `msg:"request_id"`
	Method    string        `msg:"method"`
	Args      []interface{} `msg:"args"`
}

func TypeStr(t MessageType) string {
	switch t {
	case 1:
		return "TypeRequest"
	case 2:
		return "TypeResponse"
	case 3:
		return "TypePublicEvent"
	case 4:
		return "TypePrivateEvent"
	case 5:
		return "TypeError"
	}
	return strconv.Itoa(int(t))
}
