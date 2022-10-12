package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-validator/internal/validator"
	"github.com/google/uuid"
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

	err := validator.Main().Handle(ctx, event)
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
	operationMessage := `{
	  "client_id": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e",
	  "operation": "BUY",
	  "symbol": "BTC",
	  "start_time": "2022-09-17T12:05:07.45066-03:00"
	}`

	snsEventMessage, _ := json.Marshal(createSNSEvent(operationMessage))

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
