package mocks

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type tokenBuilder struct {
	ExpectedToken string
	BuildCounter  int
	BuildError    error
}

func TokenBuilder() *tokenBuilder {
	return &tokenBuilder{}
}

func (t *tokenBuilder) Build(string, string, any, string) (string, custom_error.BaseErrorAdapter) {
	t.BuildCounter++

	if t.BuildError != nil {
		return "", exceptions.TokenBuilderError(t.BuildError, "token build error")
	}

	return t.ExpectedToken, nil
}

func (t *tokenBuilder) Reset() {
	t.ExpectedToken = ""
	t.BuildCounter = 0
	t.BuildError = nil
}
