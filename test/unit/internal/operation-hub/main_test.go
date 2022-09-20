package operation_hub

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMainSuccess(t *testing.T) {
	main := operation_hub.Main()

	assert.NotNilf(t, main, "main cannot be nil")
}
