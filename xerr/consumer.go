package xerr

type RespMsg struct {
	Ignore bool
	Error  error
}

func New() *RespMsg {
	return &RespMsg{}
}

func (r *RespMsg) Ignored() *RespMsg {
	r.Ignore = true
	return r
}

func (r *RespMsg) Err(err error) *RespMsg {
	r.Error = err
	return r
}

func (r *RespMsg) IsError() bool {
	return r.Error != nil
}

func (r *RespMsg) ShouldRequeue() bool {
	return r.Error != nil && !r.Ignore
}

func (r *RespMsg) ShouldAck() bool {
	return r.Error == nil || r.Ignore
}

func (r *RespMsg) Reset() {
	r.Ignore = false
	r.Error = nil
}
