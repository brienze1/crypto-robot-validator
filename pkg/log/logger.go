package log

import (
	"encoding/json"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

type (
	logger struct {
		Timestamp     string      `json:"timestamp"`
		Level         string      `json:"level"`
		ErrorMsg      string      `json:"exceptions,omitempty"`
		Message       string      `json:"message"`
		Metadata      interface{} `json:"metadata,omitempty"`
		CorrelationId string      `json:"correlationId"`
		TransactionId string      `json:"transactionId"`
	}
)

const (
	infoLevel  string = "INFO "
	errorLevel string = "ERROR"
)

var once sync.Once

var loggerInstance *logger

func Logger() *logger {
	if loggerInstance == nil {
		once.Do(
			func() {
				log.SetFlags(0)
				loggerInstance = &logger{}
				loggerInstance.TransactionId = uuid.New().String()
			})
	}

	return loggerInstance
}

func (l *logger) SetCorrelationID(correlationId string) {
	l.CorrelationId = correlationId
}

func (l *logger) Info(message string, metadata ...interface{}) {
	logMessage := l.generateLogMessage(infoLevel, message, nil, metadata)
	log.Println(logMessage)
}

func (l *logger) Error(err error, message string, metadata ...interface{}) {
	logMessage := l.generateLogMessage(errorLevel, message, err, metadata)
	log.Println(logMessage)
}

func (l *logger) generateLogMessage(level string, message string, err error, metadata ...[]interface{}) string {
	logg := l.clone()
	logg.Level = level
	logg.Message = message
	logg.Timestamp = time.Now().Format("2022-01-01 13:01:01")

	if len(metadata[0]) > 0 {
		logg.Metadata = metadata
	}

	if err != nil {
		errorMsg, _ := json.Marshal(custom_error.NewBaseError(err))
		logg.ErrorMsg = string(errorMsg)
	}

	logMessage, _ := json.Marshal(logg)
	return string(logMessage)
}

func (l *logger) clone() logger {
	return *l
}
