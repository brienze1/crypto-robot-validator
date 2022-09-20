package config

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/config"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
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
