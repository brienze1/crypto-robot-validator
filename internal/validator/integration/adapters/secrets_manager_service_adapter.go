package adapters

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

// SecretsManagerServiceAdapter is an adapter for secret manager service implementation.
type SecretsManagerServiceAdapter interface {
	GetSecret(secretName string, secretObject any) custom_error.BaseErrorAdapter
}
