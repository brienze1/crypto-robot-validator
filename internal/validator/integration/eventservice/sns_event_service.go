package eventservice

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	adapters2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/exceptions"
)

type snsEventService struct {
	logger adapters2.LoggerAdapter
	sns    adapters.SNSAdapter
}

func SNSEventService(logger adapters2.LoggerAdapter, sns adapters.SNSAdapter) *snsEventService {
	return &snsEventService{
		logger: logger,
		sns:    sns,
	}
}

func (s *snsEventService) Send(messageObject interface{}) error {
	s.logger.Info("Send started", messageObject)

	stringMessage, err := json.Marshal(messageObject)
	if err != nil {
		return s.abort(err, "Error while trying create string message", messageObject)
	}

	payload := string(stringMessage)
	publishInput := &sns.PublishInput{
		Message:  &payload,
		TopicArn: &properties.Properties().CryptoOperationTriggerTopicArn,
	}

	result, err := s.sns.Publish(context.TODO(), publishInput)
	if err != nil {
		return s.abort(err, "Error while trying to publish", publishInput)
	}

	s.logger.Info("Send finished", messageObject, result)
	return nil
}

func (s *snsEventService) abort(err error, message string, metadata ...interface{}) error {
	binanceWebServiceError := exceptions.SNSEventServiceError(err, message)
	s.logger.Error(binanceWebServiceError, "Send failed: "+message, metadata)
	return binanceWebServiceError
}
