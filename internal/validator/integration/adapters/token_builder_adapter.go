package adapters

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

type TokenBuilderAdapter interface {
	Build(apiSecret string, endpoint string, payload any, nonce string) (string, custom_error.BaseErrorAdapter)
}
