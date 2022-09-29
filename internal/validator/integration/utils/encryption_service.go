package utils

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
	"io"
)

type encryptionService struct {
	logger adapters.LoggerAdapter
}

func EncryptionService(logger adapters.LoggerAdapter) *encryptionService {
	return &encryptionService{
		logger: logger,
	}
}

func (e *encryptionService) AESDecrypt(hexEncryptedString string, secret string) (string, custom_error.BaseErrorAdapter) {
	e.logger.Info("AESDecrypt started")

	ciphertext, err := hex.DecodeString(hexEncryptedString)
	if err != nil {
		return "", e.abort(err, "Error while trying to decode hex string")
	}

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", e.abort(err, "Could not create new cipher")
	}

	if len(ciphertext) < aes.BlockSize {
		return "", e.abort(nil, "Text is too short")
	}

	iv := ciphertext[:aes.BlockSize]

	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	e.logger.Info("AESDecrypt finished")
	return string(ciphertext), nil
}

func (e *encryptionService) AESEncrypt(decryptedString string, secret string) (string, custom_error.BaseErrorAdapter) {
	e.logger.Info("AESEncrypt started")

	plaintext := []byte(decryptedString)

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", e.abort(err, "Could not create new cipher")
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", e.abort(err, "Error filling iv with random values")
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	hexEncryptedString := hex.EncodeToString(ciphertext)

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
