package adapters

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"

type ClientPersistenceAdapter interface {
	GetClients(config *model.ClientSearchConfig) (*[]model.Client, error)
}
