package config

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadEnvSuccess(t *testing.T) {
	err := os.Setenv("VALIDATOR_ENV", "test")
	assert.Nil(t, err)

	panicFunction := func() { config.LoadEnv() }

	assert.NotPanicsf(t, panicFunction, "Should not panic")
}

func TestLoadEnvFailure(t *testing.T) {
	err := os.Setenv("VALIDATOR_ENV", uuid.NewString())
	assert.Nil(t, err)

	panicFunction := func() { config.LoadEnv() }

	assert.Panicsf(t, panicFunction, "Should panic")
}
