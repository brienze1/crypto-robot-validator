package validator

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMainSuccess(t *testing.T) {
	main := validator.Main()

	assert.NotNilf(t, main, "main cannot be nil")
}
