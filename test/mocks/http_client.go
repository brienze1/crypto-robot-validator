package mocks

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"net/http"
	"net/http/httptest"
)

type httpClient struct {
	DoCounter          int
	DoError            error
	Server             *httptest.Server
	StatusCode         int
	GetCryptoResponse  string
	GetBalanceResponse string
}

func HttpClient() *httpClient {
	return &httpClient{
		StatusCode: 200,
	}
}

func (h *httpClient) Do(req *http.Request) (*http.Response, error) {
	h.DoCounter++

	if h.DoError != nil {
		return nil, h.DoError
	}

	realClient := http.Client{}
	return realClient.Do(req)
}

func (h *httpClient) SetupServer() {
	h.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(h.StatusCode)
		var response []byte
		if r.URL.Path == properties.Properties().BiscointGetBalancePath {
			response = []byte(h.GetBalanceResponse)
		} else if r.URL.Path == properties.Properties().BiscointGetCryptoPath {
			response = []byte(h.GetCryptoResponse)
		}

		_, _ = w.Write(response)
	}))
}

func (h *httpClient) GetUrl() string {
	return h.Server.URL
}

func (h *httpClient) Close() {
	h.Server.Close()
}

func (h *httpClient) Reset() {
	h.Server = nil
	h.DoCounter = 0
	h.DoError = nil
	h.StatusCode = 200
	h.GetCryptoResponse = ""
	h.GetBalanceResponse = ""
}
