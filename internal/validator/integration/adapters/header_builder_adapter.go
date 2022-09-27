package adapters

import "net/http"

type HeaderBuilderAdapter interface {
	BinanceHeader(clientId string) http.Header
}
