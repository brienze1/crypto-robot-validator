package utils

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/utils"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	tokenBuilder        adapters.TokenBuilderAdapter
	loggerTB            = mocks.Logger()
	encryptionServiceTB = mocks.EncryptionService()
)

var (
	apiSecret     = "9y$B?E(H+MbQeThWmZq4t7w!z%C*F)J@"
	endpoint      = "v1/balance"
	nonce         = "1234567890"
	payload       any
	expectedToken = "e69e25ede54ce43adb0ae3cd9b359ee3c6dd9301ada188155b14c27b4ce61a0f39ddc1d6be21dbe8a77ae6dd4b069e7b"
)

func setupTB() {
	loggerTB.Reset()
	encryptionServiceTB.Reset()

	tokenBuilder = utils.TokenBuilder(loggerTB, encryptionServiceTB)

	payload = map[string]string{}
}

func TestBuildSuccess(t *testing.T) {
	setupTB()

	token, err := tokenBuilder.Build(apiSecret, endpoint, payload, nonce)

	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, expectedToken, token)
	assert.Equal(t, 1, encryptionServiceTB.SHA384EncryptCounter)
	assert.Equal(t, 2, loggerTB.InfoCallCounter)
	assert.Equal(t, 0, loggerTB.ErrorCallCounter)
}

func TestBuildMarshalFailure(t *testing.T) {
	setupTB()

	payload = make(chan int)

	token, err := tokenBuilder.Build(apiSecret, endpoint, payload, nonce)

	assert.NotNil(t, err)
	assert.Equal(t, "json: unsupported type: chan int", err.Error())
	assert.Equal(t, "Payload marshal failed", err.InternalError())
	assert.Equal(t, "Error while building token", err.Description())
	assert.Equal(t, "", token)
	assert.Equal(t, 0, encryptionServiceTB.SHA384EncryptCounter)
	assert.Equal(t, 1, loggerTB.InfoCallCounter)
	assert.Equal(t, 1, loggerTB.ErrorCallCounter)
}
