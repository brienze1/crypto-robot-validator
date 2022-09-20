package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/exceptions"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
)

type handler struct {
	operationUseCase adapters.OperationUseCaseAdapter
	logger           adapters.LoggerAdapter
}

// Handler constructor method, used to inject dependencies.
func Handler(operationUseCase adapters.OperationUseCaseAdapter, logger adapters.LoggerAdapter) *handler {
	return &handler{
		operationUseCase: operationUseCase,
		logger:           logger,
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

	analysisDto := &dto.AnalysisDto{}
	if err := json.Unmarshal([]byte(snsMessage.Message), analysisDto); err != nil {
		return h.abort(err, "Error while trying to parse the analysis object")
	}

	if err := h.operationUseCase.TriggerOperations(analysisDto.Summary); err != nil {
		return h.abort(err, "Error while trying to run OperationUseCase")
	}

	h.logger.Info("Event succeeded", event, ctx)
	return nil
}

func (h *handler) abort(err error, message string) custom_error.BaseErrorAdapter {
	handlerError := exceptions.HandlerError(err, message)
	h.logger.Error(handlerError, "Event failed: "+message)
	return handlerError
}
