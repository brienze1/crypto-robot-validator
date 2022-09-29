package mocks

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
	"net/http"
)

type headerBuilder struct {
	BiscointHeaderCounter int
	BiscointHeaderError   error
}

func HeaderBuilder() *headerBuilder {
	return &headerBuilder{}
}

func (h *headerBuilder) BiscointHeader(_ string, _ string, _ any) (http.Header, custom_error.BaseErrorAdapter) {
	h.BiscointHeaderCounter++

	if h.BiscointHeaderError != nil {
		return nil, exceptions.HeaderBuilderError(h.BiscointHeaderError, "header builder error")
	}

	return http.Header{}, nil
}

func (h *headerBuilder) Reset() {
	h.BiscointHeaderCounter = 0
	h.BiscointHeaderError = nil
}
