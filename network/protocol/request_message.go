package netproto

import "github.com/google/uuid"

type RequestId = uuid.UUID

type RequestMessage interface {
	Message

	Id() RequestId
	SetId(RequestId)
}
