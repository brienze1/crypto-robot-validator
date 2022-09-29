package adapters

import (
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
	"net/http"
)

type HeaderBuilderAdapter interface {
	BiscointHeader(clientId string, endpoint string, payload any) (http.Header, custom_error.BaseErrorAdapter)
}
