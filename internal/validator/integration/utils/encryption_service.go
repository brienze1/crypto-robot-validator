package utils

import (
	"crypto"
	"crypto/aes"
	"crypto/hmac"
	"encoding/base64"
	"encoding/hex"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type encryptionService struct {
	logger adapters.LoggerAdapter
}

func EncryptionService() *encryptionService {
	return &encryptionService{}
}

func (e *encryptionService) AESDecrypt(hexEncryptedString string, secret string) (string, custom_error.BaseErrorAdapter) {
	e.logger.Info("AESDecrypt started")

	encryptedString, err := hex.DecodeString(hexEncryptedString)
	if err != nil {
		return "", e.abort(err, "Error while trying to decode hex string")
	}

	cipher, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", e.abort(err, "Error while trying to create aes cipher")
	}

	decryptedString := make([]byte, len(encryptedString))
	cipher.Decrypt(decryptedString, encryptedString)

	e.logger.Info("AESDecrypt finished")
	return string(decryptedString), nil
}

func (e *encryptionService) AESEncrypt(decryptedString string, secret string) (string, custom_error.BaseErrorAdapter) {
	e.logger.Info("AESEncrypt started")

	cipher, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", e.abort(err, "Error while trying to create cipher")
	}

	encryptedString := make([]byte, len(decryptedString))

	cipher.Encrypt(encryptedString, []byte(decryptedString))

	hexEncryptedString := hex.EncodeToString(encryptedString)

	e.logger.Info("AESEncrypt finished")
	return hexEncryptedString, nil
}

func (e *encryptionService) SHA384Encrypt(string string, secret string) string {
	e.logger.Info("SHA384Encrypt started")

	base64String := base64.StdEncoding.EncodeToString([]byte(string))

	digester := hmac.New(crypto.SHA384.New, []byte(secret))

	digester.Write([]byte(base64String))

	e.logger.Info("SHA384Encrypt finished")
	return hex.EncodeToString(digester.Sum(nil))
}

func (e *encryptionService) abort(err error, message string, metadata ...interface{}) custom_error.BaseErrorAdapter {
	encryptionServiceError := exceptions.EncryptionServiceError(err, message)
	e.logger.Error(encryptionServiceError, "Encryption service failed: "+message, metadata)
	return encryptionServiceError
}
