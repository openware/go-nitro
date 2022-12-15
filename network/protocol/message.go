package netproto

type Message interface {
	Type() string
}
