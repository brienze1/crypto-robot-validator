package adapters

import "github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"

// SecretsManagerServiceAdapter is an adapter for secret manager service implementation.
type SecretsManagerServiceAdapter interface {
	GetSecret(secretName string) (*dto.RedisSecrets, error)
}
