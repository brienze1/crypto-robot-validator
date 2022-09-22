package mocks

import "github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"

type validationUseCaseMock struct {
	ValidateCallCounter int
	ValidateError       error
}

func ValidationUseCase() *validationUseCaseMock {
	return &validationUseCaseMock{}
}

func (v *validationUseCaseMock) Validate(*model.OperationRequest) error {
	v.ValidateCallCounter++
	return v.ValidateError
}

func (v *validationUseCaseMock) Reset() {
	v.ValidateCallCounter = 0
	v.ValidateError = nil
}
