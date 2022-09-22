package config

import (
	adapters3 "github.com/brienze1/crypto-robot-validator/internal/validator/delivery/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/delivery/handler"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/usecase"
	adapters2 "github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/eventservice"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/persistence"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/webservice"
	"github.com/brienze1/crypto-robot-validator/pkg/log"
	"github.com/brienze1/crypto-robot-validator/pkg/time_utils"
	"net/http"
	"sync"
	"time"
)

var dependencyInjectorInit sync.Once
var injector *dependencyInjector

type dependencyInjector struct {
	Logger               adapters.LoggerAdapter
	HTTPClient           adapters2.HTTPClientAdapter
	DynamoDBClient       adapters2.DynamoDBAdapter
	SNSClient            adapters2.SNSAdapter
	TimeSource           adapters.TimeAdapter
	CryptoService        adapters.CryptoServiceAdapter
	ClientService        adapters.ClientServiceAdapter
	EventService         adapters.EventServiceAdapter
	ClientPersistence    adapters.ClientPersistenceAdapter
	OperationPersistence adapters.OperationPersistenceAdapter
	LockPersistence      adapters.LockPersistenceAdapter
	ValidationUseCase    adapters.ValidationUseCaseAdapter
	Handler              adapters3.HandlerAdapter
}

// DependencyInjector constructor method.
func DependencyInjector() *dependencyInjector {
	if injector == nil {
		dependencyInjectorInit.Do(func() {
			injector = &dependencyInjector{}
		})
	}

	return injector
}

// WireDependencies is used to wire the dependencies together. Also instantiates new variables in case of nil values.
func (d *dependencyInjector) WireDependencies() *dependencyInjector {
	if d.Logger == nil {
		d.Logger = log.Logger()
	}
	if d.HTTPClient == nil {
		d.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	if d.DynamoDBClient == nil {
		d.DynamoDBClient = DynamoDBClient()
	}
	if d.SNSClient == nil {
		d.SNSClient = SNSClient()
	}
	if d.TimeSource == nil {
		d.TimeSource = time_utils.Time()
	}
	if d.CryptoService == nil {
		d.CryptoService = webservice.BiscointWebService(d.Logger, d.HTTPClient)
	}
	if d.ClientService == nil {
		d.ClientService = webservice.BiscointWebService(d.Logger, d.HTTPClient)
	}
	if d.EventService == nil {
		d.EventService = eventservice.SNSEventService(d.Logger, d.SNSClient)
	}
	if d.ClientPersistence == nil {
		d.ClientPersistence = persistence.DynamoDBClientPersistence(
			d.Logger,
			d.DynamoDBClient,
		)
	}
	if d.OperationPersistence == nil {
		d.OperationPersistence = persistence.DynamoDBOperationPersistence(
			d.Logger,
			d.DynamoDBClient,
		)
	}
	if d.LockPersistence == nil {
		d.LockPersistence = persistence.RedisPersistence()
	}
	if d.ValidationUseCase == nil {
		d.ValidationUseCase = usecase.ValidationUseCase(
			d.LockPersistence,
			d.ClientPersistence,
			d.ClientService,
			d.CryptoService,
			d.OperationPersistence,
			d.EventService,
			d.Logger,
		)
	}
	if d.Handler == nil {
		d.Handler = handler.Handler(d.ValidationUseCase, d.Logger)
	}

	return d
}
