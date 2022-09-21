package config

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	config.LoadTestEnv()

	properties.Properties().Aws.Config.OverrideConfig = false
}

func TestAwsConfigSuccess(t *testing.T) {
	setup()

	snsClient := config.SNSClient()

	assert.NotNilf(t, snsClient, "Should not be nil")
}
