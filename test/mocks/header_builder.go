package mocks

import "net/http"

type headerBuilder struct {
	BinanceHeaderCounter int
}

func HeaderBuilder() *headerBuilder {
	return &headerBuilder{}
}

func (h *headerBuilder) BinanceHeader(string) http.Header {
	h.BinanceHeaderCounter++
	return http.Header{}
}

func (h *headerBuilder) Reset() {
	h.BinanceHeaderCounter = 0
}
