package mocks

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type dynamoDBClientPersistence struct {
	GetClientCounter int
	GetClientError   error
	LockCounter      int
	LockError        error
	UnlockCounter    int
	UnlockError      error
	clientsAvailable []*model.Client
}

func DynamoDBClientPersistence() *dynamoDBClientPersistence {
	return &dynamoDBClientPersistence{}
}

func (d *dynamoDBClientPersistence) GetClient(clientId string) (*model.Client, custom_error.BaseErrorAdapter) {
	d.GetClientCounter++

	if d.GetClientError != nil {
		baseError := exceptions.DynamoDBClientPersistenceError(d.GetClientError, "GetClient error")
		baseError.SetLocks(true, false)
		return nil, baseError
	}

	for _, client := range d.clientsAvailable {
		if clientId == client.Id {
			return client, nil
		}
	}
	return nil, nil
}

func (d *dynamoDBClientPersistence) Lock(client *model.Client) custom_error.BaseErrorAdapter {
	d.LockCounter++

	if d.LockError != nil {
		baseError := exceptions.DynamoDBClientPersistenceError(d.LockError, "Lock error")
		baseError.SetLocks(true, false)
		return baseError
	}

	client.Lock()

	return nil
}

func (d *dynamoDBClientPersistence) Unlock(client *model.Client) custom_error.BaseErrorAdapter {
	d.UnlockCounter++

	if d.UnlockError != nil {
		baseError := exceptions.DynamoDBClientPersistenceError(d.UnlockError, "Unlock error")
		baseError.SetLocks(true, true)
		return baseError
	}

	client.Unlock()

	return nil
}

func (d *dynamoDBClientPersistence) AddClient(client *model.Client) {
	d.clientsAvailable = append(d.clientsAvailable, client)
}

func (d *dynamoDBClientPersistence) Reset() {
	d.GetClientCounter = 0
	d.GetClientError = nil
	d.LockCounter = 0
	d.LockError = nil
	d.UnlockCounter = 0
	d.UnlockError = nil
	d.clientsAvailable = []*model.Client{}
}
