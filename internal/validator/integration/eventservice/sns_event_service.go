package eventservice

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	adapters2 "github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type snsEventService struct {
	logger adapters2.LoggerAdapter
	sns    adapters.SNSAdapter
}

// SNSEventService constructor for class.
func SNSEventService(logger adapters2.LoggerAdapter, sns adapters.SNSAdapter) *snsEventService {
	return &snsEventService{
		logger: logger,
		sns:    sns,
	}
}

// Send will create a publish request for AWS SNS topic.
func (s *snsEventService) Send(messageObject interface{}) custom_error.BaseErrorAdapter {
	s.logger.Info("Send started", messageObject)

	stringMessage, err := json.Marshal(messageObject)
	if err != nil {
		return s.abort(err, "Error while trying create string message", messageObject)
	}

	payload := string(stringMessage)
	publishInput := &sns.PublishInput{
		Message:  &payload,
		TopicArn: &properties.Properties().CryptoOperationExecutorTopicArn,
	}

	result, err := s.sns.Publish(context.TODO(), publishInput)
	if err != nil {
		return s.abort(err, "Error while trying to publish", publishInput)
	}

	s.logger.Info("Send finished", messageObject, result)
	return nil
}

func (s *snsEventService) abort(err error, message string, metadata ...interface{}) custom_error.BaseErrorAdapter {
	binanceWebServiceError := exceptions.SNSEventServiceError(err, message)
	s.logger.Error(binanceWebServiceError, "Send failed: "+message, metadata)
	return binanceWebServiceError
}
