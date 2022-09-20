package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
	"github.com/google/uuid"
	"time"
)

type ctx struct {
	context.Context
	awsRequestId string
}

func (ctx ctx) Value(any) any {
	return &lambdacontext.LambdaContext{
		AwsRequestID: ctx.awsRequestId,
	}
}

func main() {
	ctx := createContext()
	event := createSQSEvent()

	err := operation_hub.Main().Handle(ctx, event)
	if err != nil {
		panic(err)
	}
}

func createContext() *ctx {
	return &ctx{
		awsRequestId: uuid.NewString(),
	}
}

func createSQSEvent() events.SQSEvent {
	analysisDto := dto.AnalysisDto{
		Summary:   summary.StrongBuy,
		Timestamp: time.Now().Format("2022-01-01 13:01:01"),
	}

	analysisMessage, _ := json.Marshal(analysisDto)

	snsEventMessage, _ := json.Marshal(createSNSEvent(string(analysisMessage)))

	return events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(snsEventMessage),
			},
		},
	}
}

func createSNSEvent(message string) events.SNSEntity {
	return events.SNSEntity{
		Message: message,
	}
}
