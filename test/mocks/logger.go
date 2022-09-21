package mocks

type loggerMock struct {
	CorrelationId    string
	InfoCallCounter  int
	ErrorCallCounter int
}

func Logger() *loggerMock {
	return &loggerMock{}
}

func (l *loggerMock) SetCorrelationID(id string) {
	l.CorrelationId = id
}

func (l *loggerMock) Info(string, ...interface{}) {
	l.InfoCallCounter++
}

func (l *loggerMock) Error(error, string, ...interface{}) {
	l.ErrorCallCounter++
}

func (l *loggerMock) Reset() {
	l.CorrelationId = ""
	l.InfoCallCounter = 0
	l.ErrorCallCounter = 0
}
