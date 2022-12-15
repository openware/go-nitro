package netproto

type ErrorMessage struct {
	Message Message
	Error   *Error
}

var _ Message = (*ErrorMessage)(nil)

const ErrorMessageType = "error"

func (e *ErrorMessage) Type() string {
	return ErrorMessageType
}
