package mocks

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type snsClient struct {
	NumberOfMessagesSent int
}

func SNSClient() *snsClient {
	return &snsClient{}
}

func (s *snsClient) Publish(_ context.Context, _ *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	s.NumberOfMessagesSent++
	return nil, nil
}

func (s *snsClient) Reset() {
	s.NumberOfMessagesSent = 0
}
