package netproto

type ResponseMessage interface {
	Message

	RequestId() RequestId
	SetRequestId(RequestId)
}
