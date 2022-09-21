package mocks

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type dynamoDBOperationPersistence struct {
	SaveCounter         int
	SaveError           error
	operationsAvailable []*model.Operation
}

func DynamoDBOperationPersistence() *dynamoDBOperationPersistence {
	return &dynamoDBOperationPersistence{}
}

func (d *dynamoDBOperationPersistence) Save(operation *model.Operation) custom_error.BaseErrorAdapter {
	d.SaveCounter++

	if d.SaveError != nil || operation == nil {
		return exceptions.DynamoDBOperationPersistenceError(d.SaveError, "save error")
	}

	for _, operationSaved := range d.operationsAvailable {
		if operation.Id == operationSaved.Id {
			operationSaved = operation
			return nil
		}
	}

	d.operationsAvailable = append(d.operationsAvailable, operation)

	return nil
}

func (d *dynamoDBOperationPersistence) GetAllOperations() []*model.Operation {
	return d.operationsAvailable
}

func (d *dynamoDBOperationPersistence) Reset() {
	d.SaveCounter = 0
	d.SaveError = nil
	d.operationsAvailable = []*model.Operation{}
}
