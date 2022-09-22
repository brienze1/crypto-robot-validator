package adapters

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type ClientPersistenceAdapter interface {
	// GetClient will find model.Client on client repository using clientId as key
	GetClient(clientId string) (*model.Client, custom_error.BaseErrorAdapter)

	// Lock will update model.Client setting flag locked as true on client repository
	Lock(client *model.Client) custom_error.BaseErrorAdapter

	// Unlock will update model.Client setting flag locked as false on client repository
	Unlock(client *model.Client) custom_error.BaseErrorAdapter
}
