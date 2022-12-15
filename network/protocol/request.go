package netproto

type Request struct {
	id RequestId
}

func (r *Request) Id() RequestId {
	return r.id
}

func (r *Request) SetId(rid RequestId) {
	r.id = rid
}
