package actioncable

import "net/http"

type ConsumerOptions struct {
	header *http.Header
}

func NewConsumerOptions() *ConsumerOptions {
	return &ConsumerOptions{}
}

func (o *ConsumerOptions) Header() *http.Header {
	return o.header
}

func (o *ConsumerOptions) SetHeader(h *http.Header) {
	o.header = h
}
