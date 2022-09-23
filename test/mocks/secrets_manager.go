package mocks

import "github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"

type secretsManagerService struct {
	secrets map[string]*dto.RedisSecrets
}

func SecretsManagerService() *secretsManagerService {
	return &secretsManagerService{
		secrets: map[string]*dto.RedisSecrets{},
	}
}

func (s *secretsManagerService) GetSecret(secretName string) (*dto.RedisSecrets, error) {
	return s.secrets[secretName], nil
}

func (s *secretsManagerService) SetSecret(key string, value *dto.RedisSecrets) {
	s.secrets[key] = value
}
