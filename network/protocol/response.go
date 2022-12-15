package netproto

type Response struct {
	requestId RequestId
}

func (r *Response) RequestId() RequestId {
	return r.requestId
}

func (r *Response) SetRequestId(rid RequestId) {
	r.requestId = rid
}
