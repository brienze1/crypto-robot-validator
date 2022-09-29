package mocks

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type dynamoDBCredentialPersistence struct {
	GetCredentialsCounter int
	GetCredentialsError   error
	credentialsAvailable  []*dto.Credentials
}

func DynamoDBCredentialPersistence() *dynamoDBCredentialPersistence {
	return &dynamoDBCredentialPersistence{}
}

func (d *dynamoDBCredentialPersistence) GetCredentials(clientId string) (*dto.Credentials, custom_error.BaseErrorAdapter) {
	d.GetCredentialsCounter++

	if d.GetCredentialsError != nil {
		baseError := exceptions.DynamoDBCredentialsPersistenceError(d.GetCredentialsError, "GetCredentials error")
		baseError.SetLocks(true, false)
		return nil, baseError
	}

	for _, credentials := range d.credentialsAvailable {
		if clientId == credentials.ClientId {
			return credentials, nil
		}
	}
	return nil, nil
}

func (d *dynamoDBCredentialPersistence) AddCredential(credentials *dto.Credentials) {
	d.credentialsAvailable = append(d.credentialsAvailable, credentials)
}

func (d *dynamoDBCredentialPersistence) Reset() {
	d.GetCredentialsCounter = 0
	d.GetCredentialsError = nil
	d.credentialsAvailable = []*dto.Credentials{}
}
