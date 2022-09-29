package utils

import (
	"encoding/json"
	adapters2 "github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type tokenBuilder struct {
	logger            adapters2.LoggerAdapter
	encryptionService adapters.EncryptionServiceAdapter
}

func TokenBuilder(logger adapters2.LoggerAdapter, encryptionService adapters.EncryptionServiceAdapter) *tokenBuilder {
	return &tokenBuilder{
		logger:            logger,
		encryptionService: encryptionService,
	}
}

func (t *tokenBuilder) Build(apiSecret string, endpoint string, payload any, nonce string) (string, custom_error.BaseErrorAdapter) {
	t.logger.Info("Build started", endpoint, payload, nonce)

	payloadString, err := json.Marshal(payload)
	if err != nil {
		return "", t.abort(err, "Payload marshal failed")
	}

	strToBeSigned := endpoint + nonce + string(payloadString)

	t.logger.Info("Build finished", endpoint, payload, nonce, payloadString, strToBeSigned)
	return t.encryptionService.SHA384Encrypt(strToBeSigned, apiSecret), nil
}

func (t *tokenBuilder) abort(err error, message string, metadata ...interface{}) custom_error.BaseErrorAdapter {
	tokenBuilderError := exceptions.TokenBuilderError(err, message)
	t.logger.Error(tokenBuilderError, "Token builder failed: "+message, metadata)
	return tokenBuilderError
}
