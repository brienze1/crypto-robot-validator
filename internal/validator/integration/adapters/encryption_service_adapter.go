package adapters

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

type EncryptionServiceAdapter interface {
	AESDecrypt(hexEncryptedString string, secret string) (string, custom_error.BaseErrorAdapter)
	AESEncrypt(decryptedString string, secret string) (string, custom_error.BaseErrorAdapter)
	SHA384Encrypt(string string, secret string) string
}
