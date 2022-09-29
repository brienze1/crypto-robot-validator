package adapters

import "github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"

type CredentialsPersistenceAdapter interface {
	GetCredentials(id string) (dto.Credentials, error)
}
