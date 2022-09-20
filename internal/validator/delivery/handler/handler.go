package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
	"github.com/brienze1/crypto-robot-validator/internal/validator/delivery/dto"
	"github.com/brienze1/crypto-robot-validator/internal/validator/delivery/exceptions"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
)

type handler struct {
	validationUseCase adapters.ValidationUseCaseAdapter
	logger            adapters.LoggerAdapter
}

// Handler constructor method, used to inject dependencies.
func Handler(validationUseCase adapters.ValidationUseCaseAdapter, logger adapters.LoggerAdapter) *handler {
	return &handler{
		validationUseCase: validationUseCase,
		logger:            logger,
	}
}

func (h *handler) Handle(context context.Context, event events.SQSEvent) error {
	ctx, _ := lambdacontext.FromContext(context)
	h.logger.SetCorrelationID(ctx.AwsRequestID)
	h.logger.Info("Event received", event, ctx)

	snsMessage := &events.SNSEntity{}
	if err := json.Unmarshal([]byte(event.Records[0].Body), snsMessage); err != nil {
		return h.abort(err, "Error while trying to parse the SNS message")
	}

	operationRequestDto := &dto.OperationRequestDto{}
	if err := json.Unmarshal([]byte(snsMessage.Message), operationRequestDto); err != nil {
		return h.abort(err, "Error while trying to parse the operation request object")
	}

	if err := h.validationUseCase.Validate(operationRequestDto.ToModel()); err != nil {
		return h.abort(err, "Error while trying to run ValidationUseCase")
	}

	h.logger.Info("Event succeeded", event, ctx)
	return nil
}

func (h *handler) abort(err error, message string) custom_error.BaseErrorAdapter {
	handlerError := exceptions.HandlerError(err, message)
	h.logger.Error(handlerError, "Event failed: "+message)
	return handlerError
}
