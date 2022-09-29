package mocks

import (
	"encoding/json"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type secretsManagerService struct {
	GetSecretCounter int
	GetSecretError   error
	secrets          map[string][]byte
}

func SecretsManagerService() *secretsManagerService {
	return &secretsManagerService{
		secrets: map[string][]byte{},
	}
}

func (s *secretsManagerService) GetSecret(secretName string, secretObject any) custom_error.BaseErrorAdapter {
	s.GetSecretCounter++
	if s.GetSecretError != nil {
		return exceptions.SecretsManagerError(s.GetSecretError, "secrets manager error")
	}

	err := json.Unmarshal(s.secrets[secretName], secretObject)
	if err != nil {
		return exceptions.SecretsManagerError(err, "secrets manager error")
	}
	return nil
}

func (s *secretsManagerService) SetSecret(key string, value any) {
	byteValue, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	s.secrets[key] = byteValue
}

func (s *secretsManagerService) Reset() {
	s.GetSecretCounter = 0
	s.GetSecretError = nil
	s.secrets = map[string][]byte{}
}
