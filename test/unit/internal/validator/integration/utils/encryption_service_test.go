package utils

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/utils"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	encryptionService       adapters.EncryptionServiceAdapter
	loggerEncryptionService = mocks.Logger()
)

var (
	decryptedString       = "String to be encrypted"
	aesEncryptedString    = "7571a562734fe71e042c9839ea0beb80f033779aa7379b5b93f0fcc10b0a74db99712329ba92"
	sha384EncryptedString = "365af8ef6b006d2c4c8155fa4c43a0a463af0b0940a2cf0dd0159edf0f41138cd375f4b828c6ff3d64497360e5a39ac9"
	key                   = "9y$B?E(H+MbQeThWmZq4t7w!z%C*F)J@"
)

func setupEncryptionService() {
	loggerEncryptionService.Reset()

	encryptionService = utils.EncryptionService(loggerEncryptionService)
}

func TestAESDecryptSuccess(t *testing.T) {
	setupEncryptionService()

	newDecryptedString, err := encryptionService.AESDecrypt(aesEncryptedString, key)

	assert.Nil(t, err)
	assert.NotNil(t, newDecryptedString)
	assert.Equal(t, decryptedString, newDecryptedString)
}

func TestAESDecryptHexDecodeFailure(t *testing.T) {
	setupEncryptionService()

	newDecryptedString, err := encryptionService.AESDecrypt(uuid.NewString(), key)

	assert.NotNil(t, err)
	assert.Equal(t, "encoding/hex: invalid byte: U+002D '-'", err.Error())
	assert.Equal(t, "Error while trying to decode hex string", err.InternalError())
	assert.Equal(t, "Error while performing encryption", err.Description())
	assert.Equal(t, "", newDecryptedString)
}

func TestAESDecryptKeyTooBigFailure(t *testing.T) {
	setupEncryptionService()

	newDecryptedString, err := encryptionService.AESDecrypt(aesEncryptedString, uuid.NewString())

	assert.NotNil(t, err)
	assert.Equal(t, "crypto/aes: invalid key size 36", err.Error())
	assert.Equal(t, "Could not create new cipher", err.InternalError())
	assert.Equal(t, "Error while performing encryption", err.Description())
	assert.Equal(t, "", newDecryptedString)
}

func TestAESDecryptTextTooShortFailure(t *testing.T) {
	setupEncryptionService()

	newDecryptedString, err := encryptionService.AESDecrypt("61", key)

	assert.NotNil(t, err)
	assert.Equal(t, "Text is too short", err.Error())
	assert.Equal(t, "Text is too short", err.InternalError())
	assert.Equal(t, "Error while performing encryption", err.Description())
	assert.Equal(t, "", newDecryptedString)
}

func TestAESEncryptSuccess(t *testing.T) {
	setupEncryptionService()

	newEncryptedString, err := encryptionService.AESEncrypt(decryptedString, key)

	assert.Nil(t, err)
	assert.NotNil(t, newEncryptedString)

	newDecryptedString, err := encryptionService.AESDecrypt(newEncryptedString, key)

	assert.Nil(t, err)
	assert.NotNil(t, newDecryptedString)
	assert.Equal(t, decryptedString, newDecryptedString)
}

func TestAESEncryptKeyTooBigFailure(t *testing.T) {
	setupEncryptionService()

	newEncryptedString, err := encryptionService.AESEncrypt(decryptedString, uuid.NewString())

	assert.NotNil(t, err)
	assert.Equal(t, "crypto/aes: invalid key size 36", err.Error())
	assert.Equal(t, "Could not create new cipher", err.InternalError())
	assert.Equal(t, "Error while performing encryption", err.Description())
	assert.Equal(t, "", newEncryptedString)
}

func TestSHA384EncryptSuccess(t *testing.T) {
	setupEncryptionService()

	newEncryptedString := encryptionService.SHA384Encrypt(decryptedString, key)

	assert.NotNil(t, newEncryptedString)
	assert.Equal(t, sha384EncryptedString, newEncryptedString)
}
