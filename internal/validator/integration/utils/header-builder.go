package utils

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"net/http"
)

type headerBuilder struct {
	credentialsPersistence   adapters.CredentialsPersistenceAdapter
	secretsManagerService    adapters.SecretsManagerServiceAdapter
	encryptionServiceAdapter adapters.EncryptionServiceAdapter
}

func HeaderBuilder() *headerBuilder {
	return &headerBuilder{}
}

func (h *headerBuilder) BinanceHeader(clientId string) http.Header {
	return http.Header{}
}
