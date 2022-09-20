package properties

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/config"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setup() {
	os.Clearenv()
	config.LoadTestEnv()
}

func TestPropertiesMinimumCryptoSellOperationFailure(t *testing.T) {
	setup()

	_ = os.Setenv("MINIMUM_CRYPTO_SELL_OPERATION", uuid.NewString())

	panicFunction := func() { properties.Properties() }

	assert.Panicsf(t, panicFunction, "Should panic")
}

func TestPropertiesMinimumCryptoBuyOperationFailure(t *testing.T) {
	setup()

	_ = os.Setenv("MINIMUM_CRYPTO_BUY_OPERATION", uuid.NewString())

	panicFunction := func() { properties.Properties() }

	assert.Panicsf(t, panicFunction, "Should panic")
}

func TestPropertiesSuccess(t *testing.T) {
	setup()

	panicFunction := func() { properties.Properties() }

	assert.NotPanicsf(t, panicFunction, "Should not panic")
}
