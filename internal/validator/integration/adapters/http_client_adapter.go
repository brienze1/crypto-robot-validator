package adapters

import "net/http"

type HTTPClientAdapter interface {
	Do(req *http.Request) (*http.Response, error)
}
