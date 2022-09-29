package mocks

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/utils"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type encryptionService struct {
	service              adapters.EncryptionServiceAdapter
	SHA384EncryptCounter int
	AESDecryptCounter    int
	AESDecryptError      error
	AESEncryptCounter    int
}

func EncryptionService() *encryptionService {
	return &encryptionService{
		service: utils.EncryptionService(Logger()),
	}
}

func (e *encryptionService) AESDecrypt(hexEncryptedString string, secret string) (string, custom_error.BaseErrorAdapter) {
	e.AESDecryptCounter++
	if e.AESDecryptError != nil {
		return "", exceptions.EncryptionServiceError(e.AESDecryptError, "AES decrypt error")
	}
	return e.service.AESDecrypt(hexEncryptedString, secret)
}

func (e *encryptionService) AESEncrypt(decryptedString string, secret string) (string, custom_error.BaseErrorAdapter) {
	e.AESEncryptCounter++
	return e.service.AESEncrypt(decryptedString, secret)
}

func (e *encryptionService) SHA384Encrypt(string string, secret string) string {
	e.SHA384EncryptCounter++
	return e.service.SHA384Encrypt(string, secret)
}

func (e *encryptionService) Reset() {
	e.SHA384EncryptCounter = 0
	e.AESDecryptCounter = 0
	e.AESDecryptError = nil
	e.AESEncryptCounter = 0

}
