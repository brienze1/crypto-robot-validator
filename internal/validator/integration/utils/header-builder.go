package utils

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	adapters2 "github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
	"github.com/brienze1/crypto-robot-validator/pkg/time_utils"
	"net/http"
)

type headerBuilder struct {
	logger                 adapters2.LoggerAdapter
	credentialsPersistence adapters.CredentialsPersistenceAdapter
	secretsManagerService  adapters.SecretsManagerServiceAdapter
	encryptionService      adapters.EncryptionServiceAdapter
	tokenBuilder           adapters.TokenBuilderAdapter
}

func HeaderBuilder() *headerBuilder {
	return &headerBuilder{}
}

func (h *headerBuilder) BinanceHeader(clientId string, payload any) (http.Header, custom_error.BaseErrorAdapter) {
	credentials, err := h.credentialsPersistence.GetCredentials(clientId)
	if err != nil {
		return nil, h.abort(err, "Error while getting client credentials")
	}

	encryptionSecrets := &dto.EncryptionSecrets{}
	err = h.secretsManagerService.GetSecret(properties.Properties().Aws.SecretsManager.CacheSecretName, encryptionSecrets)
	if err != nil {
		return nil, h.abort(err, "Error while getting encryption key")
	}

	decryptedSecret, err := h.encryptionService.AESDecrypt(credentials.ApiSecret, encryptionSecrets.EncryptionKey)
	if err != nil {
		return nil, h.abort(err, "Error while trying to decrypt secret")
	}

	nonce := time_utils.Epoch()
	token, err := h.tokenBuilder.Build(credentials.ApiKey, decryptedSecret, payload, nonce)
	if err != nil {
		return nil, h.abort(err, "Error while trying to generate token")
	}

	return http.Header{
		"Content-Type": {"application/json"},
		"BSCNT-NONCE":  {nonce},
		"BSCNT-APIKEY": {credentials.ApiKey},
		"BSCNT-SIGN":   {token},
	}, nil
}

func (b *headerBuilder) abort(err error, message string, metadata ...interface{}) custom_error.BaseErrorAdapter {
	headerBuilderError := exceptions.HeaderBuilderError(err, message)
	b.logger.Error(headerBuilderError, "Biscoint API failed: "+message, metadata)
	return headerBuilderError
}
