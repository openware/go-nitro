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

type Message struct {
	Type      MessageType
	RequestId uint64
	Method    string
	Args      interface{}
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
