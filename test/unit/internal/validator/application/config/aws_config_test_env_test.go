package config

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupTest() {
	properties.Properties().Aws.Config.OverrideConfig = true
}

func TestAwsConfigWithOverrideFailure(t *testing.T) {
	setupTest()

	snsClient := config.SNSClient()

	assert.NotNilf(t, snsClient, "Should not be nil")
}

func TestAwsConfigWithOverrideSuccess(t *testing.T) {
	setupTest()

	endpoint, err := config.NewEndpointResolver().ResolveEndpoint("sns", "any")

	assert.Nilf(t, err, "Should be nil")
	assert.NotNilf(t, endpoint, "Should not be nil")
}
