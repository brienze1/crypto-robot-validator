package mocks

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type snsEventService struct {
	SendCounter int
	SendError   error
}

func SnsEventService() *snsEventService {
	return &snsEventService{}
}

func (s *snsEventService) Send(object interface{}) custom_error.BaseErrorAdapter {
	s.SendCounter++
	if s.SendError != nil {
		return exceptions.SNSEventServiceError(s.SendError, "Send error")
	}
	return nil
}

func (s *snsEventService) Reset() {
	s.SendCounter = 0
	s.SendError = nil
}
