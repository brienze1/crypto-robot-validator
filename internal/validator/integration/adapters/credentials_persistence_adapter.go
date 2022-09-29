package adapters

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type CredentialsPersistenceAdapter interface {
	GetCredentials(id string) (*dto.Credentials, custom_error.BaseErrorAdapter)
}
