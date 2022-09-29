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

func HeaderBuilder(
	logger adapters2.LoggerAdapter,
	credentialsPersistence adapters.CredentialsPersistenceAdapter,
	secretsManagerService adapters.SecretsManagerServiceAdapter,
	encryptionService adapters.EncryptionServiceAdapter,
	tokenBuilder adapters.TokenBuilderAdapter,
) *headerBuilder {
	return &headerBuilder{
		logger:                 logger,
		credentialsPersistence: credentialsPersistence,
		secretsManagerService:  secretsManagerService,
		encryptionService:      encryptionService,
		tokenBuilder:           tokenBuilder,
	}
}

func (h *headerBuilder) BiscointHeader(clientId string, endpoint string, payload any) (http.Header, custom_error.BaseErrorAdapter) {
	h.logger.Info("BiscointHeader started", clientId, endpoint, payload)

	credentials, err := h.credentialsPersistence.GetCredentials(clientId)
	if err != nil {
		return nil, h.abort(err, "Error while getting client credentials")
	}

	encryptionSecrets := &dto.EncryptionSecrets{}
	err = h.secretsManagerService.GetSecret(properties.Properties().Aws.SecretsManager.EncryptionSecretName, encryptionSecrets)
	if err != nil {
		return nil, h.abort(err, "Error while getting encryption key")
	}

	decryptedSecret, err := h.encryptionService.AESDecrypt(credentials.ApiSecret, encryptionSecrets.EncryptionKey)
	if err != nil {
		return nil, h.abort(err, "Error while trying to decrypt secret")
	}

	nonce := time_utils.Epoch()
	token, err := h.tokenBuilder.Build(decryptedSecret, endpoint, payload, nonce)
	if err != nil {
		return nil, h.abort(err, "Error while trying to generate token")
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("BSCNT-NONCE", nonce)
	headers.Set("BSCNT-APIKEY", credentials.ApiKey)
	headers.Set("BSCNT-SIGN", token)

	h.logger.Info("BiscointHeader finished", clientId, endpoint, payload)
	return headers, nil
}

func (h *headerBuilder) abort(err error, message string, metadata ...interface{}) custom_error.BaseErrorAdapter {
	headerBuilderError := exceptions.HeaderBuilderError(err, message)
	h.logger.Error(headerBuilderError, "Header builder failed: "+message, metadata)
	return headerBuilderError
}
